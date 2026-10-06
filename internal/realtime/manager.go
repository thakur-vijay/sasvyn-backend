package realtime

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
)

type Manager struct {
	mu          sync.RWMutex
	connections map[string]map[*websocket.Conn]struct{}
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]map[*websocket.Conn]struct{}),
	}
}

func (m *Manager) Add(userID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connections[userID] == nil {
		m.connections[userID] = make(map[*websocket.Conn]struct{})
	}

	m.connections[userID][conn] = struct{}{}
}

func (m *Manager) Remove(userID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	connections, ok := m.connections[userID]
	if !ok {
		return
	}

	delete(connections, conn)

	if len(connections) == 0 {
		delete(m.connections, userID)
	}
}

func (m *Manager) Send(
	ctx context.Context,
	userID string,
	message []byte,
) error {
	m.mu.RLock()
	connections := make([]*websocket.Conn, 0, len(m.connections[userID]))

	for conn := range m.connections[userID] {
		connections = append(connections, conn)
	}

	m.mu.RUnlock()

	for _, conn := range connections {
		if err := conn.Write(ctx, websocket.MessageText, message); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) SendEvent(
	ctx context.Context,
	userID string,
	event Event,
) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return m.Send(ctx, userID, message)
}
