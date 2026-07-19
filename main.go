package main

import (
	"fmt"
	"image/color"
	"os"

	"main/level"
	"main/player"
	"main/scroll"
	"main/server"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

var TestScrolls = []scroll.ScrollInventoryStruct{
	scroll.NewInvScroll([]string{"player_fifty_close", "start_active"}, []string{"blue"}),
	scroll.NewInvScroll([]string{"start_active"}, []string{"red"}),
	scroll.NewInvScroll([]string{"start_active"}, []string{"purple"}),
}

var Player = player.NewPlayer(utils.Vec2{X: 0, Y: 0}, TestScrolls)
var Level = level.NewLevel()

var Decided = false

func (g *Game) Update() error {
	Player.Update(&Level)
	server.Update(&Level, &Player)

	mx, my := ebiten.CursorPosition()
	utils.Mouse_X = float64(mx)
	utils.Mouse_Y = float64(my)

	if !Decided {
		if ebiten.IsKeyPressed(ebiten.KeyH) {
			go server.StartServer(&Player)
			Decided = true
			fmt.Println("Hosting server...")
		} else if ebiten.IsKeyPressed(ebiten.KeyC) {
			if server.ServerAddress == "" {
				server.ServerAddress = "localhost:8080"
			}
			go server.ConnectToServer(&Player)
			Decided = true
			fmt.Println("Connecting to:", server.ServerAddress)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{230, 230, 230, 255})

	server.Draw(screen)
	Player.Draw(screen)
	Level.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (sw, sh int) {
	return 540, 320
}

func main() {
	if len(os.Args) > 1 {
		server.ServerAddress = os.Args[1]
	}

	Player.Scrolls = TestScrolls
	Player.Health = 100

	ebiten.SetWindowSize(540*3, 320*3)
	ebiten.SetWindowTitle("Fraudsters")

	fmt.Println("Game started! Press 'H' to Host, or 'C' to Connect.")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
