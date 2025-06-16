package handlers

import (
	"encrypted-chat-relay/models"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsMessage struct {
	Content string
	SentAt  time.Time
}

type Client struct {
	UserID  uuid.UUID
	Conn    *websocket.Conn
	Handler *wsHandler
	Send    chan wsMessage
}

type wsHandler struct {
	db      *gorm.DB
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

func NewWSHandler(db *gorm.DB) *wsHandler {
	return &wsHandler{
		db:      db,
		clients: make(map[uuid.UUID]*Client),
	}
}

func (h *wsHandler) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}

	// Create a new client
	client := &Client{
		Conn:    conn,
		Handler: h,
		Send:    make(chan wsMessage, 256),
	}

	defer func() {
		h.removeClient(client)
		conn.Close()
		close(client.Send)
	}()

	fmt.Printf("New connection from: %s\n", conn.RemoteAddr().String())

	if !h.authenticate(client) {
		return
	}

	time.Sleep(time.Millisecond)
	conn.WriteMessage(websocket.TextMessage, []byte("hi from the server"))
	go h.readLoop(client)
	go h.writeLoop(client)
}

func (h *wsHandler) authenticate(client *Client) bool {
	pgp := crypto.PGP()
	var username string
	_, msg, err := client.Conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading authentication message: %v\n", err)
		client.Conn.Close()
		return false
	}
	username = string(msg)

	var user models.User
	if h.db.First(&user, "user_name = ?", username).Error != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("user not found"))
		client.Conn.Close()
		return false
	}

	chellange := uuid.NewString()

	if client.Conn.WriteMessage(websocket.TextMessage, []byte(chellange)) != nil {
		client.Conn.Close()
		return false
	}

	_, msg, err = client.Conn.ReadMessage()
	if err != nil {
		client.Conn.Close()
		return false
	}
	publicKey, err := crypto.NewKeyFromArmored(user.PublicKey)
	if err != nil {
		client.Conn.Close()
		return false
	}
	verifier, err := pgp.Verify().VerificationKey(publicKey).New()
	if err != nil {
		client.Conn.Close()
		return false
	}
	verifyResult, err := verifier.VerifyCleartext(msg)
	if err != nil {
		client.Conn.Close()
		return false
	}

	if sigErr := verifyResult.SignatureError(); sigErr != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte("invalid signature"))
		client.Conn.Close()
		return false
	}

	client.UserID = user.ID
	h.mu.Lock()
	h.clients[client.UserID] = client
	h.mu.Unlock()

	if err := client.Conn.WriteMessage(websocket.TextMessage, []byte("successfully signed in")); err != nil {
		client.Conn.Close()
		return false
	}

	return true
}

func (h *wsHandler) readLoop(client *Client) {
	defer client.Conn.Close()

	for {
		message, _, err := client.Conn.ReadMessage()
		if err != nil {
			return
		}
		fmt.Println(message)
	}
}

func (h *wsHandler) writeLoop(client *Client) {
	defer func() {
		close(client.Send)
		client.Conn.Close()
	}()

	for {
		message, ok := <-client.Send
		if !ok {
			return
		}

		err := client.Conn.WriteJSON(message)
		if err != nil {
			fmt.Printf("Error writing to websocket: %v\n", err)
			return
		}
	}
}

func (h *wsHandler) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client.UserID)
}
