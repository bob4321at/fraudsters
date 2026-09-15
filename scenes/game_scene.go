package scenes

import (
	"fmt"
	"image/color"
	"log"
	"main/level"
	"main/player"
	"main/scroll"
	"main/server"
	"main/shader"
	"main/utils"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameSceneStruct struct {
	BaseScene
	ShaderLayer *textures.Texture
	Player      player.PlayerStruct
	Level       level.LevelStruct
	Decided     bool
}

var GameScene = GameSceneStruct{}

var ChosenScrolls = []scroll.ScrollInventoryStruct{
	scroll.NewInvScroll([]string{}, []string{}),
	scroll.NewInvScroll([]string{}, []string{}),
	scroll.NewInvScroll([]string{}, []string{}),
}

func (scene *GameSceneStruct) Setup() {
	scene.SetupState = true

	scene.ShaderLayer = textures.NewTexture("./art/empty.png", shader.BorderShader)

	scene.Player = player.NewPlayer(utils.Vec2{X: 0, Y: 0}, ChosenScrolls)

	scene.Level = level.NewLevel()

	scene.Decided = false

	server.PlayerRef = &scene.Player
}

func (scene *GameSceneStruct) Draw(display *ebiten.Image) {
	display.Fill(color.RGBA{230, 230, 230, 255})

	scene.ShaderLayer.Img.Clear()

	server.Draw(scene.ShaderLayer.Img)
	scene.Player.Draw(scene.ShaderLayer.Img)
	scene.Level.Draw(scene.ShaderLayer.Img)

	scene.ShaderLayer.Draw(display, &ebiten.DrawImageOptions{})
}

func (scene *GameSceneStruct) Update() {
	scene.Player.Update(&scene.Level)
	server.Update(&scene.Level, &scene.Player)

	fmt.Println(scene.Player.Health)

	if !scene.Decided {
		if ebiten.IsKeyPressed(ebiten.KeyH) {
			go server.StartServer(&scene.Player)
			scene.Decided = true
			log.Println("Hosting server...")
		} else if ebiten.IsKeyPressed(ebiten.KeyC) {
			if server.ServerAddress == "" {
				server.ServerAddress = "localhost:8080"
			}
			go server.ConnectToServer(&scene.Player)
			scene.Decided = true
			log.Println("Connecting to:", server.ServerAddress)
		}
	}
}
