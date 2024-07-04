package t1k

import (
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type ResultFlag string

const (
	ResultFlagAllowed ResultFlag = "."
	ResultFlagBlocked ResultFlag = "?"
)

func (d ResultFlag) Byte() byte {
	return d[0]
}

type DetectorResponse struct {
	Head        byte
	Body        []byte
	Delay       []byte
	ExtraHeader []byte
	ExtraBody   []byte
	Context     []byte
	Cookie      []byte
	WebLog      []byte
	BotQuery    []byte
	BotBody     []byte
	Forward     []byte
}

func (r *DetectorResponse) Allowed() bool {
	return r.Head == ResultFlagAllowed.Byte()
}

func (r *DetectorResponse) StatusCode() int {
	str := string(r.Body)
	if str == "" {
		return http.StatusForbidden
	}
	code, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("t1k convert status code failed: %v", err)
		return http.StatusForbidden
	}
	return code
}

func (r *DetectorResponse) BlockMessage() map[string]interface{} {
	return map[string]interface{}{
		"status":   r.StatusCode(),
		"success":  false,
		"message":  "blocked by Chaitin SafeLine Web Application Firewall",
		"event_id": r.EventID(),
	}
}

func (r *DetectorResponse) EventID() string {
	extra := string(r.ExtraBody)
	if extra == "" {
		return ""
	}
	// <!-- event_id: e1impksyjq0gl92le6odi0fnobi270cj -->
	re, err := regexp.Compile(`<\!--\s*event_id:\s*([a-zA-Z0-9]+)\s*-->\s*`)
	if err != nil {
		log.Printf("t1k compile regexp failed: %v", err)
		return ""
	}
	matches := re.FindStringSubmatch(extra)
	if len(matches) < 2 {
		log.Printf("t1k regexp not match event id: %s", extra)
		return ""
	}
	return matches[1]
}

type HttpDetector struct {
	extra *HttpExtra
	req   Request
	resp  Response
}

func NewHttpDetector(req *http.Request, extra *HttpExtra) *HttpDetector {
	return &HttpDetector{
		req:   NewHttpRequest(req, extra),
		extra: extra,
	}
}

func (d *HttpDetector) SetResponse(resp *http.Response) *HttpDetector {
	d.resp = NewHttpResponse(d.req, resp, d.extra)
	return d
}

func (d *HttpDetector) DetectRequest(socket io.ReadWriter) (*DetectorResponse, error) {
	raw, err := d.req.Serialize()
	if err != nil {
		return nil, err
	}
	_, err = socket.Write(raw)
	if err != nil {
		return nil, err
	}
	return d.ReadResponse(socket)
}

func (d *HttpDetector) DetectResponse(socket io.ReadWriter) (*DetectorResponse, error) {
	raw, err := d.resp.Serialize()
	if err != nil {
		return nil, err
	}
	_, err = socket.Write(raw)
	if err != nil {
		return nil, err
	}
	return d.ReadResponse(socket)
}

func (d *HttpDetector) ReadResponse(r io.Reader) (*DetectorResponse, error) {
	res := &DetectorResponse{}
	for {
		p, err := ReadPacket(r)
		if err != nil {
			return nil, err
		}
		switch p.Tag().Strip() {
		case TAG_HEADER:
			if len(p.PayLoad()) != 1 {
				return nil, errors.New("len(T1K_HEADER) != 1")
			}
			res.Head = p.PayLoad()[0]
		case TAG_DELAY:
			res.Delay = p.PayLoad()
		case TAG_BODY:
			res.Body = p.PayLoad()
		case TAG_EXTRA_HEADER:
			res.ExtraHeader = p.PayLoad()
		case TAG_EXTRA_BODY:
			res.ExtraBody = p.PayLoad()
		case TAG_CONTEXT:
			res.Context = p.PayLoad()
		case TAG_COOKIE:
			res.Cookie = p.PayLoad()
		case TAG_WEB_LOG:
			res.WebLog = p.PayLoad()
		case TAG_BOT_QUERY:
			res.BotQuery = p.PayLoad()
		case TAG_BOT_BODY:
			res.BotBody = p.PayLoad()
		case TAG_FORWARD:
			res.Forward = p.PayLoad()
		}
		if p.Last() {
			break
		}
	}
	return res, nil
}
