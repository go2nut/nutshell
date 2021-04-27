package im

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if oldClient, exist := h.clients[client.Uid]; exist {
				close(oldClient.send)
				delete(h.clients, client.Uid)
			}
			h.clients[client.Uid] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.Uid]; ok {
				delete(h.clients, client.Uid)
				close(client.send)
			}
		case message := <-h.broadcast:
			for uid, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, uid)
				}
			}
		}
	}
}
