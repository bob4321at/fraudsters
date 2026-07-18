package player

import (
	"main/level"
	"main/scroll"
	"main/utils"
	"math"
	"sync"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayerStruct struct {
	Pos     utils.Vec2
	Vel     utils.Vec2
	CanJump bool
	Texture *textures.Texture

	Scrolls []scroll.ScrollInventoryStruct

	SelectedSpell int

	Health int
}

var ScrollQueue sync.Map

func (player *PlayerStruct) Collisions(current_level *level.LevelStruct) {
	if current_level.Collide(utils.Vec2{X: player.Pos.X + player.Vel.X, Y: player.Pos.Y}, utils.Vec2{X: 12, Y: 16}) {
		player.Vel.X = 0
	}
	if current_level.Collide(utils.Vec2{X: player.Pos.X, Y: player.Pos.Y + player.Vel.Y}, utils.Vec2{X: 12, Y: 16}) {
		player.Vel.Y = 0
		player.CanJump = true
	} else {
		player.CanJump = false
	}
}

func (player *PlayerStruct) Update(current_level *level.LevelStruct) {
	if player.Health <= 0 {
		player.Pos.Y = 10000000
	}

	player.Vel.Y += 0.1

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.Vel.X = -2
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.Vel.X = 2
	} else {
		player.Vel.X = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && player.CanJump {
		player.Vel.Y = -3
		player.CanJump = false
	}

	if ebiten.IsKeyPressed(ebiten.Key1) {
		player.SelectedSpell = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		player.SelectedSpell = 1
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		player.SelectedSpell = 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		scroll_length := 0

		ScrollQueue.Range(func(key, value any) bool {
			scroll_length += 1
			return true
		})
		// fmt.Println(scroll_length)

		VelocityAngle := math.Atan2(utils.Mouse_Y-player.Pos.Y, utils.Mouse_X-player.Pos.X)

		scroll_to_add := player.Scrolls[player.SelectedSpell]
		scroll_to_add.Velocity.X = math.Cos(VelocityAngle) * 2.5
		scroll_to_add.Velocity.Y = math.Sin(VelocityAngle) * 2.5

		ScrollQueue.Store(scroll_length+1, scroll_to_add)
	}

	player.Collisions(current_level)

	player.Pos.X += player.Vel.X
	player.Pos.Y += player.Vel.Y
}

func (player *PlayerStruct) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(int(player.Pos.X)), float64(int(player.Pos.Y)))
	player.Texture.Draw(screen, &op)
}

func NewPlayer(Pos utils.Vec2, Scrolls []scroll.ScrollInventoryStruct) PlayerStruct {
	new_player := PlayerStruct{}

	new_player.Pos = Pos
	new_player.Texture = textures.NewTexture("./art/player.png", "")
	new_player.SelectedSpell = 0

	return new_player
}
