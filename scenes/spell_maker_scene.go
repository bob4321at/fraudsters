package scenes

import (
	"main/scroll"
	"main/utils"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var SpellPickerScene SpellMakerSceneStruct

var SCROLLMAKERBACKGROUND = textures.NewTexture("./art/scroll_maker.png", "")
var BRUSHTEXTURE = textures.NewTexture("./art/brush.png", "")

type SpellMakerSceneStruct struct {
	BaseScene
}

var CurrentlyEditing = 0

var PossibleConditions []string
var PossibleSpells []string

func (scene *SpellMakerSceneStruct) Setup() {
	scene.SetupState = true
	scene.display = ebiten.NewImage(540, 320)

	for name_of_conditions := range scroll.Conditions {
		PossibleConditions = append(PossibleConditions, name_of_conditions)
	}

	for name_of_spells := range scroll.Actions {
		PossibleSpells = append(PossibleSpells, name_of_spells)
	}
}

func (scene SpellMakerSceneStruct) Draw(display *ebiten.Image) {
	SCROLLMAKERBACKGROUND.Draw(display, &ebiten.DrawImageOptions{})

	condition_count := 0

	for _, name := range PossibleConditions {
		ebitenutil.DebugPrintAt(display, name, 271, 16+condition_count*16)
		condition_count += 1
	}

	action_count := 0

	for _, name := range PossibleSpells {
		ebitenutil.DebugPrintAt(display, name, 271, 169+action_count*16)
		action_count += 1
	}

	player_conditions_total := 0
	for _, player_conditions := range ChosenScrolls[CurrentlyEditing].Conditions {
		ebitenutil.DebugPrintAt(display, player_conditions, 40, 81+player_conditions_total*16)
		player_conditions_total += 1
	}
	player_actions_total := 0
	for _, player_actions := range ChosenScrolls[CurrentlyEditing].Action {
		ebitenutil.DebugPrintAt(display, player_actions, 40, 187+player_actions_total*16)
		player_actions_total += 1
	}

	BrushDrawOps := ebiten.DrawImageOptions{}
	if CurrentlyEditing == 0 {
		BrushDrawOps.GeoM.Translate(71, 15)
	}
	if CurrentlyEditing == 1 {
		BrushDrawOps.GeoM.Translate(119, 15)
	}
	if CurrentlyEditing == 2 {
		BrushDrawOps.GeoM.Translate(167, 15)
	}
	BRUSHTEXTURE.Draw(display, &BrushDrawOps)
}

func (scene SpellMakerSceneStruct) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		if CurrentlyEditing < 2 {
			CurrentlyEditing += 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		if CurrentlyEditing > -1 {
			CurrentlyEditing -= 1
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		CurrentScene = "Game"
	}

	for i := range ChosenScrolls[CurrentlyEditing].Action {
		if utils.Collide(utils.Vec2{X: utils.Mouse_X, Y: utils.Mouse_Y}, utils.Vec2{X: 1, Y: 1}, utils.Vec2{X: 40, Y: float64(187 + i*16)}, utils.Vec2{X: 205, Y: 16}) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				utils.RemoveArrayElement(i, &ChosenScrolls[CurrentlyEditing].Action)
				i = len(ChosenScrolls[CurrentlyEditing].Action) + 1
			}
		}
	}

	for i, condition := range PossibleConditions {
		if utils.Collide(utils.Vec2{X: utils.Mouse_X, Y: utils.Mouse_Y}, utils.Vec2{X: 1, Y: 1}, utils.Vec2{X: 271, Y: float64(16 + i*16)}, utils.Vec2{X: 1000, Y: 16}) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				can_add := true
				for _, other_conditions := range ChosenScrolls[CurrentlyEditing].Conditions {
					if other_conditions == condition {
						can_add = false
					}
				}

				if !can_add {
					break
				}

				ChosenScrolls[CurrentlyEditing].Conditions = append(ChosenScrolls[CurrentlyEditing].Conditions, condition)
			}
		}
	}

	for i, spell := range PossibleSpells {
		if utils.Collide(utils.Vec2{X: utils.Mouse_X, Y: utils.Mouse_Y}, utils.Vec2{X: 1, Y: 1}, utils.Vec2{X: 271, Y: float64(169 + i*16)}, utils.Vec2{X: 1000, Y: 16}) {
			if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
				can_add := true
				for _, other_spells := range ChosenScrolls[CurrentlyEditing].Action {
					if other_spells == spell {
						can_add = false
					}
				}

				if !can_add {
					break
				}

				ChosenScrolls[CurrentlyEditing].Action = append(ChosenScrolls[CurrentlyEditing].Action, spell)
			}
		}
	}
}

var SpellMakerScene = SpellMakerSceneStruct{}
