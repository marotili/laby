// Copyright (c) 2012 by Lecture Hall Games Authors.
// All source files are distributed under the Simplified BSD License.

package main

import (
	"encoding/gob"
	// "fmt"
	"github.com/banthar/Go-SDL/mixer"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/banthar/Go-SDL/ttf"
	// "github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"github.com/banthar/gl"
	"go/build"
	"laby/game"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

const basePkg = "github.com/fruhwirth-marco/lecture-hall-games"

type Player struct {
	Conn      net.Conn
	Nick      string
	ButtonA   bool
	ButtonB   bool
	JoystickX float32
	JoystickY float32
}

const (
	screenWidth  = 800
	screenHeight = 900
)

// type Game interface {
// 	Update(t time.Duration)
// 	Render(screen *sdl.Surface)
// 	Join(player *Player, x, y float32)
// 	Leave(player *Player)
// 	KeyPressed(input sdl.Keysym)
// }

var (
	// game Game
	mu sync.Mutex
)

func PollEvents() []sdl.Event {
	events := make([]sdl.Event, 0)
	for {
		event := sdl.PollEvent()
		if event == nil {
			return events
		}

		events = append(events, event)
	}
}

func main() {
	log.SetFlags(log.Llongfile)
	runtime.LockOSThread()

	conn, err := net.Dial("tcp", "129.27.19.194:8001")
	if err != nil {
		log.Fatal("No connection to server")
		return
	}

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	var screen = sdl.SetVideoMode(screenWidth, screenHeight, 32, sdl.OPENGL|sdl.HWSURFACE|sdl.GL_DOUBLEBUFFER)
	if screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption("Lecture Hall Games", "")
	sdl.EnableUNICODE(1)
	if gl.Init() != 0 {
		log.Fatal("could not initialize OpenGL")
	}

	gl.Viewport(0, 0, int(screen.W), int(screen.H))
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(screen.W), float64(screen.H), 0, -1.0, 1.0)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.DEPTH_TEST)
	gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)

	if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 4096) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	if p, err := build.Default.Import(basePkg, "", build.FindOnly); err == nil {
		os.Chdir(p.Dir)
	}

	// rand.Seed(time.Now().UnixNano())
	// levelDir := fmt.Sprintf("data/levels/demolevel%d", 3+rand.Intn(numberLevels))
	//carsDir := fmt.Sprintf(" data/cars/car%d/", 1+rand.Intn(numberCars))
	// if game, err = NewGame(); err != nil {
	// 	log.Fatal(err)
	// }

	running := true
	last := time.Now()

	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	var player game.Player
	dec.Decode(&player)
	log.Println("Received player self")

	renderData := LoadRenderData()
	clientGame, _ := game.NewGame()
	// get player id from server
	clientGame.NewPlayer(int(player))

	is := game.NewInputState(clientGame, player)

	gameStarted := false

	log.Println("We are", player)

	var music *mixer.Music
	var font *ttf.Font

	if music = mixer.LoadMUS("data/music.ogg"); music == nil {
		log.Fatal(sdl.GetError())
	}

	mixer.ResumeMusic()
	// music.PlayMusic(-1)

	if font = ttf.OpenFont("data/font.otf", 32); font == nil {
		log.Fatal(sdl.GetError())
	}
	textStart := ttf.RenderUTF8_Blended(font, "Hello", sdl.Color{0, 0, 0, 0})
	spriteStart := NewSpriteFromSurface(textStart)

	for running {
		Clear()
		for _, event := range PollEvents() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.ResizeEvent:
				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					if e.Keysym.Sym == sdl.K_ESCAPE {
						running = false
					} else {
						// game.KeyPressed(e.Keysym)
					}
				}

				is.HandleEvent(e)
			}
		}

		// }

		current := time.Now()
		t := current.Sub(last)
		last = current

		// handle user input
		playerActions := is.StepActions(t)

		clientGame.Update(t)

		if len(playerActions) > 1 {
			// log.Fatal("Sending multiple actions not supported")
		}

		if !gameStarted {
			// log.Println("Request game state")
			var otherPlayer game.Player
			var otherPlayerJoined bool
			var gameStartsNow bool

			enc.Encode(game.ClientReqGameState)
			dec.Decode(&otherPlayerJoined)
			if otherPlayerJoined {
				dec.Decode(&otherPlayer)
				dec.Decode(&gameStartsNow)
				clientGame.NewPlayer(int(otherPlayer))
				log.Println("Added other player")
				gameStarted = gameStartsNow
			}
		}

		// send user input to server

		// log.Println("Send new actions")
		filteredActions := make([]game.ActionType, 0)
		for _, action := range playerActions {
			enc.Encode(game.ClientReqSendAction)
			enc.Encode(1)
			enc.Encode(action)

			var serverResp game.ServerResponse
			dec.Decode(&serverResp)
			if serverResp != game.ServerActionOk {
				log.Println("server action not ok", serverResp)
				// drop actions
				// log.Println("Dropped actions")
			} else {
				filteredActions = append(filteredActions, action)
			}
		}

		// now fetch input from other users

		// log.Println("Requesting client update")
		enc.Encode(game.ClientReqUpdate)
		var numPlayers int
		var numActions int
		var otherPlayer game.Player
		var action game.ActionType
		dec.Decode(&numPlayers)

		var data map[game.Player][]game.ActionType = make(map[game.Player][]game.ActionType, 0)
		for i := 0; i < numPlayers; i++ {
			dec.Decode(&otherPlayer)
			data[otherPlayer] = make([]game.ActionType, 0)

			dec.Decode(&numActions)
			for j := 0; j < numActions; j++ {
				dec.Decode(&action)
				log.Println("Received action", otherPlayer, action)
				data[otherPlayer] = append(data[otherPlayer], action)
			}
		}

		data[player] = playerActions
		for thePlayer, actions := range data {
			for _, action := range actions {
				log.Println("Perform action from player", thePlayer, action)
				err := clientGame.PerformPlayerAction(thePlayer, action)
				log.Println(err)
			}
		}
		RenderMap(player, renderData, clientGame)

		if !gameStarted {
			spriteStart.Draw(50, 50, 0, 1, true)
		}

		// TODO
		// selfPlayer := game.Player(0)

		// UpdateGame

		sdl.GL_SwapBuffers()
	}

	sdl.Quit()
}
