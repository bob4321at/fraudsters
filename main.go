package main

import (
	"log"
	"os"

	"main/scenes"
	"main/server"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	utils.Mouse_X = float64(mx)
	utils.Mouse_Y = float64(my)

	if !scenes.SceneList[scenes.CurrentScene].GetSetup() {
		scenes.SceneList[scenes.CurrentScene].Setup()
	}

	scenes.SceneList[scenes.CurrentScene].Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if scenes.SceneList[scenes.CurrentScene].GetSetup() {
		scenes.SceneList[scenes.CurrentScene].Draw(screen)
	}
}

func (g *Game) Layout(ow, oh int) (sw, sh int) {
	return 540, 320
}

func main() {
	if len(os.Args) > 1 {
		server.ServerAddress = os.Args[1]
	}

	ebiten.SetWindowSize(540*3, 320*3)
	ebiten.SetWindowTitle("Fraudsters")

	log.Println("Game started! Press 'H' to Host, or 'C' to Connect.")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
