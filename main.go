package main

import (
	"image/color"
	"log"
	"os"

	"main/level"
	"main/player"
	"main/scroll"
	"main/server"
	"main/shader"
	"main/utils"

	"github.com/bob4321at/textures"
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

var ShaderLayer = textures.NewTexture("./art/empty.png", shader.BorderShader)

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
			log.Println("Hosting server...")
		} else if ebiten.IsKeyPressed(ebiten.KeyC) {
			if server.ServerAddress == "" {
				server.ServerAddress = "localhost:8080"
			}
			go server.ConnectToServer(&Player)
			Decided = true
			log.Println("Connecting to:", server.ServerAddress)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{230, 230, 230, 255})

	ShaderLayer.Img.Clear()

	server.Draw(ShaderLayer.Img)
	Player.Draw(ShaderLayer.Img)
	Level.Draw(ShaderLayer.Img)

	ShaderLayer.Draw(screen, &ebiten.DrawImageOptions{})
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

	log.Println("Game started! Press 'H' to Host, or 'C' to Connect.")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
