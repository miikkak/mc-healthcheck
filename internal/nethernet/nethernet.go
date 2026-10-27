// Package nethernet implements a minimal Bedrock NetherNet signaling check:
// just enough to confirm a Geyser listener is up and serving its signaling
// document. The WebRTC (UDP) gameplay leg is deliberately not exercised.
package nethernet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxResponseLength = 1 << 20 // guards against a misbehaving peer streaming an unbounded body

// Status performs GET /v1/join against the signaling endpoint at host:port
// and returns an error unless it answers 200 with a JSON object carrying a
// protocol number within timeout. useTLS selects https over http.
func Status(host string, port int, timeout time.Duration, useTLS bool) error {
	scheme := "http"
	if useTLS {
		scheme = "https"
	}
	u := url.URL{Scheme: scheme, Host: net.JoinHostPort(host, strconv.Itoa(port)), Path: "/v1/join"}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("get %s: %w", u.String(), err)
	}

	// Redirects are not followed: a 3xx from the configured endpoint must
	// not let another server's 200 pass for this listener.
	client := &http.Client{
		Transport: &http.Transport{DisableKeepAlives: true},
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("get %s: %w", u.String(), err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get %s: unexpected status %s", u.String(), resp.Status)
	}

	var status struct {
		Protocol *int `json:"protocol"`
	}
	dec := json.NewDecoder(io.LimitReader(resp.Body, maxResponseLength))
	if err := dec.Decode(&status); err != nil {
		return fmt.Errorf("get %s: malformed response: %w", u.String(), err)
	}
	// Exactly one JSON document: anything but trailing whitespace is malformed.
	if err := dec.Decode(&json.RawMessage{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("get %s: malformed response: trailing data after JSON document", u.String())
	}
	if status.Protocol == nil {
		return fmt.Errorf("get %s: response has no protocol field", u.String())
	}
	return nil
}
