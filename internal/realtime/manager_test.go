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

func TestManagerDisconnectClient_ClosesOnlyMatchingUserClientConnections(t *testing.T) {
	manager := NewManager()
	connected := make(chan struct{}, 3)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}

		userID := r.URL.Query().Get("user")
		clientID := r.Header.Get(clientIDHeader)
		manager.Add(userID, clientID, conn)
		connected <- struct{}{}
		defer manager.Remove(userID, conn)
		defer conn.Close(websocket.StatusNormalClosure, "")

		for {
			if _, _, err := conn.Read(r.Context()); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	dial := func(userID, clientID string) *websocket.Conn {
		t.Helper()

		url := strings.Replace(server.URL, "http://", "ws://", 1) + "?user=" + userID
		conn, _, err := websocket.Dial(context.Background(), url, &websocket.DialOptions{
			HTTPHeader: http.Header{clientIDHeader: []string{clientID}},
		})
		if err != nil {
			t.Fatalf("dial websocket for %s/%s: %v", userID, clientID, err)
		}
		t.Cleanup(func() {
			_ = conn.CloseNow()
		})
		return conn
	}

	disconnectedConn := dial("user-1", "device-a")
	dial("user-1", "device-b")
	dial("user-2", "device-a")

	for range 3 {
		select {
		case <-connected:
		case <-time.After(time.Second):
			t.Fatal("timed out waiting for websocket connections to register")
		}
	}

	if err := manager.DisconnectClient("user-1", "device-a"); err != nil {
		t.Fatalf("disconnect client: %v", err)
	}

	readCtx, cancelRead := context.WithTimeout(context.Background(), time.Second)
	defer cancelRead()
	if _, _, err := disconnectedConn.Read(readCtx); err == nil {
		t.Fatal("expected the matching websocket connection to close")
	}

	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if len(manager.connections["user-1"]) != 1 {
		t.Errorf("user-1 connections = %d, want 1", len(manager.connections["user-1"]))
	}
	if len(manager.connections["user-2"]) != 1 {
		t.Errorf("user-2 connections = %d, want 1", len(manager.connections["user-2"]))
	}
}
