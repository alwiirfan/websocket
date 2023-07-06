package websocket

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			if _, ok := hub.Rooms[client.RoomID]; ok {
				room := hub.Rooms[client.RoomID]

				if _, ok := room.Clients[client.ID]; !ok {
					room.Clients[client.ID] = client
				}
			}
		case client := <-hub.Unregister:
			if _, ok := hub.Rooms[client.RoomID].Clients[client.ID]; ok {
				if _, ok := hub.Rooms[client.RoomID].Clients[client.ID]; ok {
					// broadcast a message saying that the client has left the room
					if len(hub.Rooms[client.RoomID].Clients) != 0 {
						hub.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   client.RoomID,
							Username: client.Username,
						}
					}

					delete(hub.Rooms[client.RoomID].Clients, client.ID)
					close(client.Message)
				}
			}
		case message := <-hub.Broadcast:
			if _, ok := hub.Rooms[message.RoomID]; ok {

				for _, client := range hub.Rooms[message.RoomID].Clients {
					client.Message <- message
				}
			}
		}

	}
}
