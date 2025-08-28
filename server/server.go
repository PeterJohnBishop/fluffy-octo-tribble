package server

import (
	"fluffy-coto-tribble/server/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	hub := newHub()
	go hub.run()

	// WebSocket
	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})

	// connect DynamoDB
	dynamoClient := services.ConnectDB()
	AddDynamoDBRoutes(dynamoClient, router)

	// connect S3
	s3Client := services.ConnectS3()
	AddS3Routes(s3Client, dynamoClient, router)

	// connect Google Maps
	mapClient := services.FindMaps()
	AddMapRoutes(mapClient, router)

	log.Println("Server listening on :8080")
	router.Run(":8080")
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) run() {

	log.Println("WebSocket server listening on /ws")

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client connected")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("Client disconnected")
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allow all origins TESTING ONLY
		return true
	},
}

func serveWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan []byte, 256)}

	hub.register <- client

	go client.write()
	go client.read(hub)
}

func (c *Client) read(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		hub.broadcast <- msg
	}
}

func (c *Client) write() {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
