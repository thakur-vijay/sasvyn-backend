package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestPublisherPublish_SkipsOriginClient(t *testing.T) {
	manager := NewManager()
	connected := make(chan struct{}, 2)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		manager.Add("user-1", r.Header.Get(clientIDHeader), conn)
		connected <- struct{}{}
		defer manager.Remove("user-1", conn)
		defer conn.Close(websocket.StatusNormalClosure, "")

		for {
			if _, _, err := conn.Read(r.Context()); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	dial := func(clientID string) *websocket.Conn {
		t.Helper()

		conn, _, err := websocket.Dial(context.Background(), strings.Replace(server.URL, "http://", "ws://", 1), &websocket.DialOptions{
			HTTPHeader: http.Header{clientIDHeader: []string{clientID}},
		})
		if err != nil {
			t.Fatalf("dial websocket for %s: %v", clientID, err)
		}
		t.Cleanup(func() {
			_ = conn.CloseNow()
		})
		return conn
	}

	originConn := dial("device-a")
	recipientConn := dial("device-b")

	for range 2 {
		select {
		case <-connected:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for websocket connections to register")
		}
	}

	publisher := NewPublisher(manager)
	if err := publisher.Publish(
		context.Background(),
		"user-1",
		Event{Type: "test.event", Data: map[string]string{"value": "ok"}},
		"device-a",
	); err != nil {
		t.Fatalf("publish event: %v", err)
	}

	readCtx, cancelRead := context.WithTimeout(context.Background(), time.Second)
	_, message, err := recipientConn.Read(readCtx)
	cancelRead()
	if err != nil {
		t.Fatalf("read event on other client: %v", err)
	}

	var got Event
	if err := json.Unmarshal(message, &got); err != nil {
		t.Fatalf("decode received event: %v", err)
	}
	if got.Type != "test.event" {
		t.Fatalf("received event type %q, want %q", got.Type, "test.event")
	}

	originCtx, cancelOrigin := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelOrigin()
	_, _, err = originConn.Read(originCtx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("origin client read error = %v, want context deadline exceeded", err)
	}
}

func TestPublisherPublish_RequiresClientID(t *testing.T) {
	publisher := NewPublisher(NewManager())

	if err := publisher.Publish(context.Background(), "user-1", Event{Type: "test.event"}); !errors.Is(err, ErrClientIDRequired) {
		t.Fatalf("publish error = %v, want %v", err, ErrClientIDRequired)
	}
}
