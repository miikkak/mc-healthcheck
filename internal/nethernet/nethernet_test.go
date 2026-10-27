package nethernet

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func serve(t *testing.T, h http.HandlerFunc) (string, int) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	host, p, err := net.SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatalf("split addr: %v", err)
	}
	port, err := strconv.Atoi(p)
	if err != nil {
		t.Fatalf("parse port: %v", err)
	}
	return host, port
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
	}{
		{"valid", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"protocol":2225,"version":"26.60"}`))
		}, false},
		{"missing protocol", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`{"version":"26.60"}`))
		}, true},
		{"html", func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(`<html></html>`))
		}, true},
		{"not found", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var path string
			host, port := serve(t, func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				tc.handler(w, r)
			})
			err := Status(host, port, 2*time.Second, false)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Status() error = %v, wantErr %v", err, tc.wantErr)
			}
			if path != "/v1/join" {
				t.Errorf("path = %q, want /v1/join", path)
			}
		})
	}
}

func TestStatus_TimesOut(t *testing.T) {
	host, port := serve(t, func(http.ResponseWriter, *http.Request) { time.Sleep(time.Second) })
	if err := Status(host, port, 100*time.Millisecond, false); err == nil {
		t.Fatal("Status() = nil, want timeout error")
	}
}

func TestStatus_ConnectionRefused(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	if err := Status("127.0.0.1", port, time.Second, false); err == nil {
		t.Fatal("Status() = nil, want connection error")
	}
}
