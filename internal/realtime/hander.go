package realtime

import (
	"net/http"

	"github.com/coder/websocket"
	"github.com/sasvyn/backend/internal/modules/auth"
)

type Handler struct {
	manager *Manager
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}

	h.manager.Add(userID, conn)

	defer h.manager.Remove(userID, conn)
	defer conn.Close(websocket.StatusNormalClosure, "")

	for {
		_, _, err := conn.Read(r.Context())
		if err != nil {
			break
		}
	}

	<-r.Context().Done()
}
