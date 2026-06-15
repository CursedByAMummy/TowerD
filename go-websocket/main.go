package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"math"
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
    ID    string
    Money int
    Alive bool

    Base    Base
    Turrets []Turret
    Units   []Unit
}

type GameState struct {
	Players []PlayerState `json:"players"`
}

type Position struct {
	X 	float64:`json:"x"`
	Y 	float64:`json:"y"`
}

type Base struct {
    Position Position `json:"position"`
    Health   int      `json:"health"`
}

type Unit struct {
    ID       int      `json:"id"`
    Position Position `json:"position"`

    Health int `json:"health"`
    Damage int `json:"damage"`
    Speed  int `json:"speed"`

    TargetPlayer string `json:"targetPlayer"`
}

type Turret struct {
    ID int

    Position Position

    Health int
    Damage int

    Range float64
    ROF   float64
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
				X: p.X,
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

	for _, player := range Players {
		player.Conn.WriteJSON(state)
	}
}

func moveUnits(player *PlayerState, game *Game) {
	for _, p := range player.Units {
		unit := &player.Units[i]

		targetPlayer := game.Players[unit.TargetPlayer]

		moveUnit(unit, &targetPlayer.Base,
		)
	}
}

func moveUnit(unit *Unit, targetBase *Base) {
	for _, o := range Players[p][offenses] {
		target := players[p][offtarget]
		if targetBase.Position.X > unit.Position.X {
			unit.Position.X += 1
		}
		if targetBase.Position.X < unit.Position.X {
			unit.Position.X -= 1
		}
		if targetBase.Position.Y > unit.Position.Y {
			unit.Position.Y += 1
		}
		if targetBase.Position.Y < unit.Position.Y {
			unit.Position.Y -= 1
		}
	return
	}	
}

func updateCombat() {
	for _, p := range Players {
		if players[p][alive] == true {
			for _, d := range Players[p][defenses] {
				defTarget := findDefTarget(p, d)
				dealDamage(defTarget)
			}
		}
	}
}

func inRange(turret Turret, unit Unit) bool {
	dx := turret.Position.X - unit.Position.X
	dy := turret.Position.Y - unit.Position.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance <= turret.Range
}

func findDefTarget(p int, d int) target []byte {
	target := [0,0]
	minDist := 100000000
	for _, p1 := range Players {
		if players[p1][alive] == true && p != p1 {
			for _, o := range Players[p1][offenses] {
				distX := int(players[p][defenses][d][X] - players[p][offenses][o][X])
				distY := int(players[p][defenses][d][Y] - players[p][offenses][o][Y])
				hypo := int(math.Sqrt(int(math.Pow(distX, 2)) + int(math.Pow(distY, 2))))
				if minDist > hypo {
					minDist = hypo
					target = [p1, o]
				}
			}
		}
	}
	return target
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}
	if len(players) >= 4 {
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
	fmt.Printf("Received: %s\n", msg)
	if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		fmt.Println("Error writing message:", err)
		break
	}
}

func gameLoop() {
    ticker := time.NewTicker(
        100 * time.Millisecond,
    )
    for range ticker.C {
		for _, p := range Players {
			if players[p][alive] == true {
				moveUnits(p)
			}
			updateCombat()
//        collectResources()
//        checkVictory()
//        broadcastState()
    }
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	gameLoop()
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