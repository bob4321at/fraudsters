package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SCREEN_WIDTH  = 540
	SCREEN_HEIGHT = 320
)

type BaseScene struct {
	display    *ebiten.Image
	SetupState bool
}

type Scene interface {
	Setup()
	DeSetup()
	Update()
	Draw(display *ebiten.Image)
	GetSetup() bool
}

func (scene *BaseScene) Setup() {
	scene.SetupState = true
	scene.display = ebiten.NewImage(540, 320)
}

func (scene *BaseScene) DeSetup() {
	scene.SetupState = false
}

func (scene *BaseScene) Update() {}

func (scene *BaseScene) Draw(display *ebiten.Image) {}

func (scene *BaseScene) GetSetup() bool {
	return scene.SetupState
}

var SceneList = map[string]Scene{
	"Game":         &GameScene,
	"Spell Picker": &SpellPickerScene,
}

var CurrentScene = "Spell Picker"
