package server

import (
	"flag"
	"os"
	"time"

	"v/internal/handlers" // handlers for our various use cases
	w "v/pkg/webrtc"      // webrtc package imported as w

	"github.com/gofiber/fiber/v2" // go fiber framework
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2" //websocket support
)

var (
	addr = flag.String("addr", ":"+os.Getenv("PORT"), "") // address to listen on
	cert = flag.String("cert", "", "") // TLS certificate file path
	key  = flag.String("key", "", "") // TLS key file path
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080" // default to localhost and port 8080
	}

	// initialize html template engine
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine}) // create a new fiber app instance
	app.Use(logger.New()) // use logger middleware for logging purposes
	app.Use(cors.New()) // use  Cross-Origin Resource Sharing middleware for handling cors headers

	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate) // our room.go handler routed on /room/create, but replaces create by uuid which is link to share with others to join the meeting room
	app.Get("/room/:uuid", handlers.Room) // redirected to the uuid
	app.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{ // creating new websocket
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/room/:uuid/chat", handlers.RoomChat) // chat interface
	app.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket)) // gives the room chat websocket connection for sending messages
	app.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket)) 
	app.Get("/stream/:suuid", handlers.Stream) // stream connection for our video conferencing
	app.Get("/stream/:suuid/websocket", websocket.New(handlers.StreamWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/stream/:suuid/chat/websocket", websocket.New(handlers.StreamChatWebsocket))
	app.Get("/stream/:suuid/viewer/websocket", websocket.New(handlers.StreamViewerWebsocket)) // room viewr webscoket connection, video conferencing
	app.Static("/", "./assets") // load all static assests needed for the project

	// Initialize maps to store rooms and streams
	w.Rooms = make(map[string]*w.Room)
	w.Streams = make(map[string]*w.Room)
	// Start a goroutine to dispatch key frames periodically, threads in go for concurrent programming
	go dispatchKeyFrames()
	if *cert != "" {
		return app.ListenTLS(*addr, *cert, *key)
	}
	return app.Listen(*addr)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range w.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}
