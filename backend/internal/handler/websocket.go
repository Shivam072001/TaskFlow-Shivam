package handler

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/shivam/taskflow/backend/internal/dto"
	"github.com/shivam/taskflow/backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	hub       *ws.Hub
	jwtSecret string
}

func NewWebSocketHandler(hub *ws.Hub, jwtSecret string) *WebSocketHandler {
	return &WebSocketHandler{hub: hub, jwtSecret: jwtSecret}
}

func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		dto.WriteError(w, http.StatusUnauthorized, "missing token")
		return
	}

	userID, email, tokenExp, err := h.parseToken(tokenStr)
	if err != nil {
		dto.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := ws.NewClient(h.hub, conn, userID, email, tokenExp)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *WebSocketHandler) parseToken(tokenStr string) (uuid.UUID, string, time.Time, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil {
		return uuid.Nil, "", time.Time{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, "", time.Time{}, jwt.ErrTokenInvalidClaims
	}

	userIDStr, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, "", time.Time{}, err
	}

	var exp time.Time
	if expFloat, ok := claims["exp"].(float64); ok {
		exp = time.Unix(int64(expFloat), 0)
	}

	return userID, email, exp, nil
}
