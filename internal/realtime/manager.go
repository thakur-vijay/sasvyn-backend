package realtime

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/coder/websocket"
)

type Manager struct {
	mu          sync.RWMutex
	connections map[string]map[*websocket.Conn]struct{}
}

func NewManager() *Manager {
	log.Println("[Realtime] Creating WebSocket manager")

	return &Manager{
		connections: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (m *Manager) Add(userID string, conn *websocket.Conn) {
	log.Printf("[Realtime] Adding connection for user: %s", userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connections[userID] == nil {
		log.Printf("[Realtime] Creating connection pool for user: %s", userID)

		m.connections[userID] = make(map[*websocket.Conn]struct{})
	}

	m.connections[userID][conn] = struct{}{}

	log.Printf(
		"[Realtime] Connection added | user=%s | connections=%d",
		userID,
		len(m.connections[userID]),
	)
}

func (m *Manager) Remove(userID string, conn *websocket.Conn) {
	log.Printf("[Realtime] Removing connection for user: %s", userID)

	m.mu.Lock()
	defer m.mu.Unlock()

	connections, ok := m.connections[userID]
	if !ok {
		log.Printf("[Realtime] No connection pool found for user: %s", userID)
		return
	}

	delete(connections, conn)

	log.Printf(
		"[Realtime] Connection removed | user=%s | remaining=%d",
		userID,
		len(connections),
	)

	if len(connections) == 0 {
		delete(m.connections, userID)

		log.Printf("[Realtime] Removed empty connection pool for user: %s", userID)
	}
}

func (m *Manager) Send(
	ctx context.Context,
	userID string,
	message []byte,
) error {
	log.Printf(
		"[Realtime] Sending message | user=%s | payload_size=%d bytes",
		userID,
		len(message),
	)

	m.mu.RLock()

	connections := make([]*websocket.Conn, 0, len(m.connections[userID]))

	for conn := range m.connections[userID] {
		connections = append(connections, conn)
	}

	connectionCount := len(connections)

	m.mu.RUnlock()

	log.Printf(
		"[Realtime] Broadcasting message | user=%s | connections=%d",
		userID,
		connectionCount,
	)

	var sendErr error

	for _, conn := range connections {
		if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
			log.Printf(
				"[Realtime] Failed to send message | user=%s | error=%v",
				userID,
				err,
			)

			// Remove stale/dead connection.
			m.Remove(userID, conn)

			// Don't stop broadcasting to other devices.
			sendErr = err
			continue
		}

		log.Printf(
			"[Realtime] Message sent successfully | user=%s",
			userID,
		)
	}

	return sendErr
}

func (m *Manager) SendEvent(
	ctx context.Context,
	userID string,
	event Event,
) error {
	log.Printf(
		"[Realtime] Sending event | user=%s | event=%+v",
		userID,
		event,
	)

	message, err := json.Marshal(event)
	if err != nil {
		log.Printf(
			"[Realtime] Failed to encode event | user=%s | error=%v",
			userID,
			err,
		)

		return err
	}

	log.Printf(
		"[Realtime] Event encoded | user=%s | payload_size=%d bytes",
		userID,
		len(message),
	)

	return m.Send(ctx, userID, message)
}
