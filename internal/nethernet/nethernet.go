// Package nethernet implements a minimal Bedrock NetherNet signaling check:
// just enough to confirm a Geyser listener is up and serving its signaling
// document. The WebRTC (UDP) gameplay leg is deliberately not exercised.
package nethernet

import (
	"encoding/json"
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

	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{DisableKeepAlives: true},
	}
	resp, err := client.Get(u.String())
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
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseLength)).Decode(&status); err != nil {
		return fmt.Errorf("get %s: malformed response: %w", u.String(), err)
	}
	if status.Protocol == nil {
		return fmt.Errorf("get %s: response has no protocol field", u.String())
	}
	return nil
}
