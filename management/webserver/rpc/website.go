package rpc

import (
	"encoding/json"
	"io"
	"sync"
	"time"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	pb "chaitin.cn/patronus/safeline-2/management/webserver/proto/website"
)

var (
	ping = pb.Event{Type: Ping, Msg: nil}

	// subscriberMu guards subscriber, which is written by the goroutine that
	// serves the Subscribe stream and read by every request handler that
	// publishes a configuration.
	subscriberMu sync.RWMutex
	subscriber   *StreamClient // only ONE sub
)

func Publish(msg []byte, eventType string) error {
	// The subscriber is read once and every use below goes through the local
	// variable: a new stream may replace it at any moment, and re-reading the
	// global pointer would let this call talk to a closing stream or to a
	// different subscription than the one it started with.
	sub := currentSubscriber()
	if sub == nil {
		return errors.New("Service is abnormal, and the nginx conf cannot be updated for the time being. Please go to the shell to check the relevant logs")
	}

	// Only one exchange may be outstanding on the stream. The response carries
	// no request id, so two overlapping callers that share the single slot below
	// would each read the answer of the other: one would commit a database
	// change that tcontrollerd refused, the other would roll back one it
	// applied. The generated client does not allow concurrent sends on one
	// stream either.
	sub.exchangeMu.Lock()
	defer sub.exchangeMu.Unlock()

	// A response whose own request already gave up is still in the slot, and
	// this call would read it as its own answer.
	discardPendingResponses(sub)

	err := sub.stream.Send(&pb.Event{
		Type: eventType,
		Msg:  msg,
	})
	if err != nil {
		select {
		case <-sub.stream.Context().Done():
			err = sub.stream.Context().Err() // context canceled
			clearSubscriber(sub)
			return err
		case sub.errCh <- errors.Wrapf(err, "Send event err: %s", msg):
			return errors.Wrapf(err, "Send event err: %s", msg)
		}
	}

	rspTimeoutTicker := time.NewTicker(WaitRspTimeout)
	defer rspTimeoutTicker.Stop()

	for {
		select {
		case rsp := <-sub.rspCh:
			if rsp.Err {
				return errors.New(string(rsp.Msg))
			} else {
				// success
				return nil
			}
		case <-rspTimeoutTicker.C:
			// The device did not answer, so what it applied is unknown: an
			// answer that arrives later would be read by the next request as
			// if it belonged to it. The stream is dropped instead, which
			// leaves the slot this call abandons unreachable and makes
			// tcontrollerd reconnect and synchronise with a full push.
			clearSubscriber(sub)
			select {
			case sub.errCh <- errors.New("Wait timeout for updating result"):
			default:
			}

			return errors.New("Wait timeout for updating result")
		}
	}
}

// discardPendingResponses drops the responses that are sitting in the slot.
//
// It is called while the exchange of the caller is held, so nothing in the slot
// can belong to a request that is still waiting.
func discardPendingResponses(sc *StreamClient) {
	for {
		select {
		case rsp := <-sc.rspCh:
			logger.Warnf("Discarding an unclaimed control channel response of type %s", rsp.GetType())
		default:
			return
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

	// exchangeMu serialises the send-and-wait of Publish. The response of a
	// request is not labelled with the request it answers, so at most one
	// request may be in flight at a time.
	exchangeMu sync.Mutex
}

// currentSubscriber returns the active subscriber, or nil when no stream is
// connected.
func currentSubscriber() *StreamClient {
	subscriberMu.RLock()
	defer subscriberMu.RUnlock()

	return subscriber
}

// setSubscriber makes sc the active subscriber.
func setSubscriber(sc *StreamClient) {
	subscriberMu.Lock()
	defer subscriberMu.Unlock()

	if subscriber != nil && subscriber != sc {
		// A client that restarts reconnects while its previous stream may not
		// have timed out yet, so a replacement is expected; it is logged
		// because only a client that knows the control channel token can
		// trigger it.
		logger.Warn("A new subscriber replaces the previous one")
	}

	subscriber = sc
}

// clearSubscriber forgets the active subscriber, but only when it is still the
// one that asks: a stream that ends after it was replaced must not disconnect
// its successor.
func clearSubscriber(sc *StreamClient) {
	subscriberMu.Lock()
	defer subscriberMu.Unlock()

	if subscriber == sc {
		subscriber = nil
	}
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
		if current := currentSubscriber(); current != nil && current != sc {
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
		if current := currentSubscriber(); current != nil && current != sc {
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

			// This loop is also the goroutine that reads the keepalive pongs
			// and resets the timer, so it must not block on a full slot: it
			// would disconnect a stream whose peer is perfectly healthy. Only
			// one exchange is outstanding at a time (see Publish) and it
			// empties the slot before it sends, so a slot that is already full
			// holds a response nobody is waiting for.
			select {
			case sc.rspCh <- rsp:
			default:
				logger.Warn("Discarding a control channel response that no request is waiting for")
			}
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
	// UnimplementedWebsiteServer is embedded through a pointer, and the code
	// that registers a service checks that the pointer is not nil: the check
	// is there because a nil pointer would only be noticed when a method that
	// this server does not implement is called. It has to be set here, or the
	// registration panics as soon as the generated code carries that check.
	return &WebsiteServer{UnimplementedWebsiteServer: &pb.UnimplementedWebsiteServer{}}
}

func publishFullWebsite() error {
	var websites []model.Website
	db := database.GetDB().DB

	// A failed query leaves the slice empty, and an empty list means "this
	// installation has no sites" to tcontrollerd: it would delete every nginx
	// configuration it manages. The error has to stop the push instead.
	if err := db.Model(&model.Website{}).Find(&websites).Error; err != nil {
		return errors.Wrap(err, "load websites")
	}

	if websites == nil {
		// A nil slice marshals to "null", which tcontrollerd reads as a
		// malformed push. An installation that really has no site sends an
		// empty list.
		websites = []model.Website{}
	}

	byteWebsites, err := json.Marshal(&websites)
	if err != nil {
		return err
	}

	return Publish(byteWebsites, EventTypeFullWebsite)
}

// Subscribe is gRPC API entrypoint
func (ws *WebsiteServer) Subscribe(stream pb.Website_SubscribeServer) error {
	sc := newStreamClient(stream)
	setSubscriber(sc)

	defer sc.timer.Stop()
	// A stream that ends without being replaced leaves the channel unusable
	// until a new client connects, so the global pointer is released here.
	defer clearSubscriber(sc)

	go sc.pingLoop()
	go sc.recvLoop()

	if err := publishFullWebsite(); err != nil {
		// triggered when tcd starts, ignore push error messages
		logger.Warn(err)
	}

	select {
	case <-sc.quit:
		logger.Infof("Disconnected gracefully")
		return nil
	case <-sc.timer.C:
		logger.Error("Keepalive timeout")
	case err := <-sc.errCh:
		logger.WithError(err).Error()
		return err
	case <-sc.stream.Context().Done():
		logger.Infof("Subscribe context done: %s", sc.stream.Context().Err())
	}

	return nil
}
