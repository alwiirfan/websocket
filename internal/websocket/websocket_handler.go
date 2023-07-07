package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (handler *Handler) CreateRoom(ctx *gin.Context) {
	var request CreateRoomRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler.hub.Rooms[request.ID] = &Room{
		ID:      request.ID,
		Name:    request.Name,
		Clients: make(map[string]*Client),
	}

	ctx.JSON(http.StatusOK, request)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000"
	},
}

func (handler *Handler) JoinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// websocket/joinRoom/:roomId?userId=1&username=user

	roomID := ctx.Param("roomId")
	clientID := ctx.Query("userId")
	username := ctx.Query("username")

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	message := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// register a new client through the register channer
	handler.hub.Register <- client

	// broadcast that message
	handler.hub.Broadcast <- message

	go client.WriteMessage()
	client.ReadMessage(handler.hub)
}

type RoomResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (handler *Handler) GetRooms(ctx *gin.Context) {
	rooms := make([]RoomResponse, 0)

	for _, room := range handler.hub.Rooms {
		rooms = append(rooms, RoomResponse{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	ctx.JSON(http.StatusOK, rooms)
}

type ClientResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (handler *Handler) GetClients(ctx *gin.Context) {
	var clients []ClientResponse
	roomId := ctx.Param("roomId")

	if _, ok := handler.hub.Rooms[roomId]; !ok {
		clients = make([]ClientResponse, 0)
		ctx.JSON(http.StatusOK, clients)
	}

	for _, c := range handler.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientResponse{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	ctx.JSON(http.StatusOK, clients)
}
