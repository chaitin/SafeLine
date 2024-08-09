package t1k

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/chaitin/t1k-go/detection"
	"github.com/chaitin/t1k-go/t1k"

	"github.com/chaitin/t1k-go/misc"
)

func writeDetectionRequest(w io.Writer, req detection.Request) error {
	{
		data, err := req.Header()
		if err != nil {
			return err
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_HEADER|t1k.MASK_FIRST, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return err
		}
	}
	{
		bodySize, bodyReadCloser, err := req.Body()
		if err == nil {
			defer bodyReadCloser.Close()
			sec := t1k.MakeReaderSection(t1k.TAG_BODY, bodySize, bodyReadCloser)
			err = t1k.WriteSection(sec, w)
			if err != nil {
				return err
			}
		}
	}
	{
		data, err := req.Extra()
		if err != nil {
			return err
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_EXTRA, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return err
		}
	}
	{
		sec := t1k.MakeSimpleSection(t1k.TAG_VERSION|t1k.MASK_LAST, []byte("Proto:2\n"))
		err := t1k.WriteSection(sec, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeDetectionResponse(w io.Writer, rsp detection.Response) error {
	{
		data, err := rsp.RequestHeader()
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_HEADER|t1k.MASK_FIRST, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
	}
	{
		data, err := rsp.Header()
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_RSP_HEADER, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
	}
	{
		bodySize, bodyReadCloser, err := rsp.Body()
		if err == nil {
			defer bodyReadCloser.Close()
			sec := t1k.MakeReaderSection(t1k.TAG_RSP_BODY, bodySize, bodyReadCloser)
			err = t1k.WriteSection(sec, w)
			if err != nil {
				return err
			}
		}
	}
	{
		data, err := rsp.Extra()
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_RSP_EXTRA, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
	}
	{
		sec := t1k.MakeSimpleSection(t1k.TAG_VERSION, []byte("Proto:2\n"))
		err := t1k.WriteSection(sec, w)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
	}
	{
		data, err := rsp.T1KContext()
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
		sec := t1k.MakeSimpleSection(t1k.TAG_CONTEXT|t1k.MASK_LAST, data)
		err = t1k.WriteSection(sec, w)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
	}
	return nil
}

func readDetectionResult(r io.Reader) (*detection.Result, error) {
	var ret detection.Result
	parseSection := func(sec t1k.Section) error {
		var buf bytes.Buffer
		err := sec.WriteBody(&buf)
		if err != nil {
			return misc.ErrorWrap(err, "")
		}
		tag := sec.Header().Tag.Strip()
		switch tag {
		case t1k.TAG_HEADER:
			if len(buf.Bytes()) != 1 {
				return fmt.Errorf("len(T1K_HEADER) != 1")
			}
			ret.Head = buf.Bytes()[0]
		case t1k.TAG_BODY:
			ret.Body = buf.Bytes()
		case t1k.TAG_ALOG:
			ret.Alog = buf.Bytes()
		case t1k.TAG_EXTRA_HEADER:
			ret.ExtraHeader = buf.Bytes()
		case t1k.TAG_EXTRA_BODY:
			ret.ExtraBody = buf.Bytes()
		case t1k.TAG_CONTEXT:
			ret.T1KContext = buf.Bytes()
		case t1k.TAG_COOKIE:
			ret.Cookie = buf.Bytes()
		case t1k.TAG_WEB_LOG:
			ret.WebLog = buf.Bytes()
		}
		return nil
	}
	sec, err := t1k.ReadFullSection(r)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	if !sec.Header().Tag.IsFirst() {
		return nil, fmt.Errorf("first section IsFirst != true, middle of another msg or corrupt stream, with <%x>", sec.Header().Tag)
	}
	for {
		err = parseSection(sec)
		if err != nil {
			return nil, misc.ErrorWrap(err, "")
		}
		if sec.Header().Tag.IsLast() {
			break
		}
		sec, err = t1k.ReadSection(r)
		if err != nil {
			return nil, misc.ErrorWrap(err, "")
		}
	}
	return &ret, nil
}

func doDetectRequest(s io.ReadWriter, req detection.Request) (*detection.Result, error) {
	err := writeDetectionRequest(s, req)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	ret, err := readDetectionResult(s)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	ret.Objective = detection.RO_REQUEST
	return ret, nil
}

func doDetectResponse(s io.ReadWriter, rsp detection.Response) (*detection.Result, error) {
	err := writeDetectionResponse(s, rsp)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	ret, err := readDetectionResult(s)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	ret.Objective = detection.RO_RESPONSE
	return ret, nil
}

func DetectRequestInCtx(s io.ReadWriter, dc *detection.DetectionContext) (*detection.Result, error) {
	ret, err := doDetectRequest(s, dc.Request)
	if err != nil {
		return nil, err
	}
	dc.ProcessResult(ret)
	return ret, nil
}

func DetectResponseInCtx(s io.ReadWriter, dc *detection.DetectionContext) (*detection.Result, error) {
	ret, err := doDetectResponse(s, dc.Response)
	if err != nil {
		return nil, misc.ErrorWrap(err, "")
	}
	dc.ProcessResult(ret)
	return ret, nil
}

func Detect(s io.ReadWriter, dc *detection.DetectionContext) (*detection.Result, *detection.Result, error) {
	var reqResult *detection.Result
	var rspResult *detection.Result
	if dc.Request != nil {
		ret, err := doDetectRequest(s, dc.Request)
		if err != nil {
			return nil, nil, misc.ErrorWrap(err, "")
		}
		reqResult = ret
		dc.ProcessResult(reqResult)
	}
	if dc.Response != nil {
		ret, err := doDetectResponse(s, dc.Response)
		if err != nil {
			return nil, nil, misc.ErrorWrap(err, "")
		}
		rspResult = ret
		dc.ProcessResult(rspResult)
	}
	return reqResult, rspResult, nil
}

func DetectHttpRequest(s io.ReadWriter, req *http.Request) (*detection.Result, error) {
	dc, _ := detection.MakeContextWithRequest(req)
	return doDetectRequest(s, detection.MakeHttpRequestInCtx(req, dc))
}

func DetectRequest(s io.ReadWriter, req detection.Request) (*detection.Result, error) {
	return doDetectRequest(s, req)
}
