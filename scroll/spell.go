package scroll

import (
	"encoding/json"
	"fmt"
	"main/utils"
	"math"
	"os"
	"sync"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
)

type NetworkedSpell struct {
	Name        string
	Position    utils.Vec2
	Velocity    utils.Vec2
	Damage      int
	Speed       int
	Lifetime    int
	Movement    string
	ImgPath     string
	Destroy     bool
	JustSpawned bool
	Size        utils.Vec2
}

type SpellInfo struct {
	HostPosistion   utils.Vec2
	ClientPosistion utils.Vec2
}

type Condition struct {
	Name  string
	Value int
}

var ConditionsToApply sync.Map

func AddSpell(path string, ScrollPos utils.Vec2) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var spell NetworkedSpell
	json.Unmarshal(bytes, &spell)

	img := textures.NewTexture(spell.ImgPath, "")

	img_width := img.Img.Bounds().Dx()
	img_height := img.Img.Bounds().Dy()

	spell.Size.X = float64(img_width)
	spell.Size.Y = float64(img_height)

	spell.Position = utils.Vec2{X: ScrollPos.X - float64(img_width/2), Y: ScrollPos.Y - float64(img_height/2)}

	spell_queue_length := 0

	SpellQueue.Range(func(key, value any) bool {
		spell_queue_length += 1
		return true
	})

	SpellQueue.Store(spell_queue_length+1, spell)
}

var SpellMovement = map[string]func(this_spell *NetworkedSpell, info *SpellInfo){
	"Static":  StaticSpell,
	"ShootAt": ShootAt,
}

func (spell *NetworkedSpell) Draw(screen *ebiten.Image) {
	texture := textures.NewTexture(spell.ImgPath, "")

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(spell.Position.X, spell.Position.Y)

	texture.Draw(screen, &op)
}

func (spell *NetworkedSpell) Update(info *SpellInfo) {
	SpellMovement[spell.Movement](spell, info)
}

func StaticSpell(current_spell *NetworkedSpell, info *SpellInfo) {}

func ShootAt(current_spell *NetworkedSpell, info *SpellInfo) {
	if !current_spell.JustSpawned {
		client_dist := utils.GetDist(current_spell.Position, info.ClientPosistion)
		host_dist := utils.GetDist(current_spell.Position, info.HostPosistion)

		fmt.Println(info.ClientPosistion)

		target_angle := 90.0

		if client_dist <= host_dist {
			target_angle = (math.Atan2(current_spell.Position.Y-info.ClientPosistion.Y, current_spell.Position.X-info.ClientPosistion.X))
		} else {
			target_angle = (math.Atan2(current_spell.Position.Y-info.HostPosistion.Y, current_spell.Position.X-info.HostPosistion.X))
		}

		current_spell.Velocity.X = -math.Cos(target_angle) * float64(current_spell.Speed)
		current_spell.Velocity.Y = -math.Sin(target_angle) * float64(current_spell.Speed)
	}

	current_spell.Position.X += current_spell.Velocity.X
	current_spell.Position.Y += current_spell.Velocity.Y

	current_spell.JustSpawned = true
}
