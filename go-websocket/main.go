package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)
	// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	players = make(map[string]*Player)
	mutex 	sync.RWMutex
)

type Game struct {
	Players map[string]*PlayerState
	Resources []Resource
}

type PlayerState struct {
	ID			string
	Money		int
	Alive		bool
	
	Home		Home
	Defenses	[]Defenses
	Offenses	[]Offenses
}

type Player struct {
	ID 			string
	Username	string
	Conn 		*websocket.Conn
}

type ClientMessage struct {
	Type		string `json:"type"`
	Direction	string `json:"direction"`
}

type PlayerState struct {
	ID string 	`json:"id"`
	X int		`json:"x"`
	Y int		`json:"y"`
}

type GameState struct {
	Players []PlayerState `json:"players"`
}

func getGameState() GameState {
	state := GameState{}

	mutex.RLock()
	defer mutex.RUnlock()

	for _, p := range players {
		state.Players = append(
			state.Players,
			PlayerState{
				ID: p.ID,
				X: p.X
				Y: p.Y,
			},
		)
	}
	return state
}

func broadcastState() {
	state := getGameState()

	mutex.RLock()
	defer mutex.RUnlock()

	for _, player := range players {
		player.Conn.WriteJSON(state)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	mutex.RLock()
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}
	if len(players) >= 4 {
		mutex.RUnlock()
		http.Error(w, "Game Full", http.StatusForbidden)
		return
	}
	id := fmt.Sprintf(
		"player%d",
		len(players)+1,
	)

	player := &Player{
		ID:		id,
		Username:	userID,
		Conn:		conn,
		X:		0,
		Y:		0,
	}
	players[id] = player
	mutex.RUnlock()

	defer conn.Close()

	go handleConnection(player)
}

func handleConnection(player *Player) {
	conn := player.Conn
	for {
		var msg ClientMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		switch msg.Type {
		case "move":
		}
	}
	fmt.Printf("Received: %s\n", message)
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		fmt.Println("Error writing message:", err)
		break
	}
}
														  }

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}



`"switch msg.Type {

case "move":

    switch msg.Direction {

    case "up":
        player.Y--

    case "down":
        player.Y++

    case "left":
        player.X--

    case "right":
        player.X++
    }

    fmt.Printf(
        "%s moved to (%d,%d)\n",
        player.ID,
        player.X,
        player.Y,
    )
}"`