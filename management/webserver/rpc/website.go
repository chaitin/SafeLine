package rpc

import (
	"encoding/json"
	"io"
	"time"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	pb "chaitin.cn/patronus/safeline-2/management/webserver/proto/website"
)

var (
	ping       = pb.Event{Type: Ping, Msg: nil}
	Subscriber *StreamClient // only ONE sub
)

func Publish(msg []byte, eventType string) error {
	if Subscriber == nil {
		return errors.New("Service is abnormal, and the nginx conf cannot be updated for the time being. Please go to the shell to check the relevant logs")
	}
	err := Subscriber.stream.Send(&pb.Event{
		Type: eventType,
		Msg:  msg,
	})
	if err != nil {
		select {
		case <-Subscriber.stream.Context().Done():
			err = Subscriber.stream.Context().Err() // context canceled
			Subscriber = nil
			return err
		case Subscriber.errCh <- errors.Wrapf(err, "Send event err: %s", msg):
			return errors.Wrapf(err, "Send event err: %s", msg)
		}
	}

	rspTimeoutTicker := time.NewTicker(WaitRspTimeout)
	defer rspTimeoutTicker.Stop()

	for {
		select {
		case rsp := <-Subscriber.rspCh:
			if rsp.Err {
				return errors.New(string(rsp.Msg))
			} else {
				// success
				return nil
			}
		case <-rspTimeoutTicker.C:
			return errors.New("Wait timeout for updating result")
		}
	}
}

// StreamClient is the instance for every client stream connected
type StreamClient struct {
	stream pb.Website_SubscribeServer
	timer  *time.Timer // timer for timeout
	errCh  chan error
	rspCh  chan *pb.Response
	quit   chan struct{} // quit stream client gracefully
}

func newStreamClient(stream pb.Website_SubscribeServer) *StreamClient {
	return &StreamClient{
		stream: stream,
		timer:  time.NewTimer(KeepaliveTimeout),
		errCh:  make(chan error, 1),
		rspCh:  make(chan *pb.Response, 1),
		quit:   make(chan struct{}),
	}
}

func (sc *StreamClient) pingLoop() {
	pingTicker := time.NewTicker(KeepaliveTime)
	defer pingTicker.Stop()

	for {
		if Subscriber != nil && sc != Subscriber {
			// new subscriber in replace of the old one.
			logger.Debug("New subscriber in replace of the old one")
			return
		}

		select {
		case <-sc.stream.Context().Done():
			return
		case <-pingTicker.C:
			err := sc.stream.Send(&ping)
			if err != nil {
				select {
				case <-sc.stream.Context().Done():
				case sc.errCh <- errors.Wrapf(err, "Send ping err"):
				}
				return
			}
		}
	}
}

func (sc *StreamClient) recvLoop() {
	for {
		if Subscriber != nil && sc != Subscriber {
			// new subscriber in replace of the old one.
			close(sc.quit)
			logger.Debug("New subscriber in replace of the old one")
			return
		}

		rsp, err := sc.stream.Recv()
		if err != nil {
			if err != io.EOF {
				select {
				case <-sc.stream.Context().Done():
				case sc.errCh <- errors.Wrapf(err, "Recv pong err"):
				}
			}
			return
		}

		logger.Debugf("Got message Type %s", rsp.GetType())
		if rsp.Type != Pong {
			// receive response
			logger.Infof("Recv updating website rsp: err(%t), msg(%s)", rsp.GetErr(), rsp.GetMsg())
			sc.rspCh <- rsp
			continue
		}
		sc.timer.Reset(KeepaliveTimeout)
	}
}

// WebsiteServer is the gRPC server implementation
type WebsiteServer struct {
	*pb.UnimplementedWebsiteServer
}

func GetWebsiteServer() *WebsiteServer {
	return &WebsiteServer{}
}

func publishFullWebsite() error {
	var websites []model.Website
	db := database.GetDB().DB
	db.Model(&model.Website{}).Find(&websites)
	byteWebsites, err := json.Marshal(&websites)
	if err != nil {
		return err
	}

	return Publish(byteWebsites, EventTypeFullWebsite)
}

// Subscribe is gRPC API entrypoint
func (ws *WebsiteServer) Subscribe(stream pb.Website_SubscribeServer) error {
	Subscriber = newStreamClient(stream)
	defer Subscriber.timer.Stop()

	go Subscriber.pingLoop()
	go Subscriber.recvLoop()

	if err := publishFullWebsite(); err != nil {
		// triggered when tcd starts, ignore push error messages
		logger.Warn(err)
	}

	select {
	case <-Subscriber.quit:
		logger.Infof("Disconnected gracefully")
		return nil
	case <-Subscriber.timer.C:
		logger.Error("Keepalive timeout")
	case err := <-Subscriber.errCh:
		logger.WithError(err).Error()
		return err
	case <-Subscriber.stream.Context().Done():
		logger.Infof("Subscribe context done: %s", Subscriber.stream.Context().Err())
		Subscriber = nil
	}

	return nil
}
