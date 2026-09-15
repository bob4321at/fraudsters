package level

import (
	"main/utils"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
)

// TODO ADD LEVELS
type LevelStruct struct {
	LevelTex *textures.Texture
}

func (level *LevelStruct) Draw(screen *ebiten.Image) {
	level.LevelTex.Draw(screen, &ebiten.DrawImageOptions{})
}

func (level *LevelStruct) Collide(o_pos, o_size utils.Vec2) bool {
	if o_pos.Y+o_size.Y > 288 {
		return true
	}

	if o_pos.X < 0 {
		return true
	}

	if o_pos.X+o_size.X > 540 {
		return true
	}

	if utils.Collide(o_pos, o_size, utils.Vec2{X: 192, Y: 253}, utils.Vec2{X: 156, Y: 5}) {
		return true
	}

	if utils.Collide(o_pos, o_size, utils.Vec2{X: 20, Y: 219}, utils.Vec2{X: 156, Y: 5}) {
		return true
	}

	if utils.Collide(o_pos, o_size, utils.Vec2{X: 364, Y: 219}, utils.Vec2{X: 156, Y: 5}) {
		return true
	}

	return false
}

func NewLevel() LevelStruct {
	new_level := LevelStruct{}

	new_level.LevelTex = textures.NewTexture("./art/temp_level.png", "")

	return new_level
}
