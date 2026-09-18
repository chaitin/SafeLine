package fvm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/dev/go/log"
)

const (
	// pushPolicyAttempts is how often a failed push is retried.
	pushPolicyAttempts = 3

	// maxPushPolicyResponseBody bounds the part of an error response that is
	// kept for the message.
	maxPushPolicyResponseBody = 4 << 10
)

// pushPolicy sends a compiled policy to the detector.
//
// The body of an http.Request is a stream that can be read once, so a retry
// that reuses a request sends an empty body. The detector accepts an empty
// update with 200, which would leave the previous policy in place while the
// caller believes the new one was applied: every attempt therefore builds its
// own request around the same bytes.
func pushPolicy(snURL string, policy []byte) error {
	client := &http.Client{}

	var lastErr error
	for attempt := 0; attempt < pushPolicyAttempts; attempt++ {
		req, err := http.NewRequest(http.MethodPost, snURL, bytes.NewReader(policy))
		if err != nil {
			return errors.Annotatef(err, "create request failed")
		}
		req.Header.Add("Content-Type", "application/octet-stream")
		req.ContentLength = int64(len(policy))

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, readErr := io.ReadAll(io.LimitReader(resp.Body, maxPushPolicyResponseBody))
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Warnf("Close push FSL response failed: %s", closeErr)
		}
		if readErr != nil {
			log.Warnf("Read push FSL response failed: %s", readErr)
		}

		if resp.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("Push FSL response(%d %s)", resp.StatusCode, string(body)))
		}

		return nil
	}

	return errors.Annotatef(lastErr, "Get response failed")
}
