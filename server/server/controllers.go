package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sas604/IOT-fun/server/plug"
)

type switchState struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Id    string `json:"id"`
}

type Client struct {
	Broker *WsBroker

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

type WsBroker struct {
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (c *Client) writePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.send {
		err := c.Conn.WriteJSON(string(message))
		if err != nil {
			log.Println(err)
		}
	}
}

func NewWsBroker() *WsBroker {
	return &WsBroker{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
func (WsB *WsBroker) Start() {
	for {
		select {
		case client := <-WsB.register:
			fmt.Println("registe", client)
			WsB.clients[client] = true
		case client := <-WsB.unregister:
			if _, ok := WsB.clients[client]; ok {
				delete(WsB.clients, client)
				close(client.send)
			}
		case message := <-WsB.Broadcast:
			fmt.Println(WsB.clients)
			for client := range WsB.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(WsB.clients, client)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{}

func ServeWs(WsB *WsBroker, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	client := &Client{Broker: WsB, Conn: conn, send: make(chan []byte)}
	client.Broker.register <- client
	go client.writePump()
}
func GetAllSwitches(c *gin.Context) {
	res, err := GetAllSwitchesWithState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")

	}
	c.IndentedJSON(200, res)
}

func SetSwitch(c *gin.Context) {
	var r switchState
	err := c.ShouldBindJSON(&r)
	if (err) != nil {
		c.JSON(http.StatusInternalServerError, "error processing post request")
	}
	s := plug.NewSwitch(r.Id)
	err = s.SetSwitchState(r.State)
	if (err) != nil {
		c.JSON(http.StatusInternalServerError, "error processing post request")
	}
	c.JSON(200, "Change Succesfull")
}
