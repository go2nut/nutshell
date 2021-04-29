package im

import (
	"context"
	"fmt"
	"log"
	relCli "nutshell/_example/apps/rel/client"
	"nutshell/_example/apps/shard"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[int64]*Client

	// Inbound messages from the clients.
	broadcast chan *Msg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

}

type Msg struct {
	user *shard.User
	body []byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Msg),
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
			log.Printf("unregister user:%d", client.Uid)
			if _, ok := h.clients[client.Uid]; ok {
				delete(h.clients, client.Uid)
				close(client.send)
			}
		case message := <-h.broadcast:
			for uid, client := range h.clients {
				if resp, err := relCli.Client.IsFriend(context.Background(), &shard.UserPairRequest{
					Uid1:                 uid,
					Uid2:                 message.user.UserId,
				}); err == nil {
					msgBody := message.body
				    if uid == message.user.UserId {
						msgBody = []byte(fmt.Sprintf("[you] %s", string(message.body)))
					} else if resp.IsFriend {
						msgBody = []byte(fmt.Sprintf("[friend] %s", string(message.body)))
					} else {
						msgBody = []byte(fmt.Sprintf("[stranger] %s", string(message.body)))
					}
					select {
					    case client.send <- msgBody:
					    default:
					    	close(client.send)
					    	delete(h.clients, uid)
					}
				} else {
					log.Printf("error check if friend uid1:%d uid2:%d", uid, message.user)
				}
			}
		}
	}
}
