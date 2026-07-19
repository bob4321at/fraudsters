package server

import (
	"main/level"
	"main/player"
	"main/scroll"
	"main/utils"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ServerSideScrolls struct {
	ID         int
	Posistion  utils.Vec2
	Conditions []string
	Action     []string
	Vel        utils.Vec2
	Trigged    bool
	Lifetime   int
}

func AddScrolls(Player *player.PlayerStruct) {
	player.ScrollQueue.Range(func(key, value any) bool {
		current_scroll_to_add_data, ok := player.ScrollQueue.Load(key)
		if ok {
			current_scroll_to_add := current_scroll_to_add_data.(scroll.ScrollInventoryStruct)

			server_side_scroll := ServerSideScrolls{}
			server_side_scroll.Action = current_scroll_to_add.Action
			server_side_scroll.Conditions = current_scroll_to_add.Conditions
			server_side_scroll.Posistion = Player.Pos
			server_side_scroll.ID = len(GameState.Scrolls) + 1
			server_side_scroll.Vel = current_scroll_to_add.Velocity

			GameState.Scrolls = append(GameState.Scrolls, server_side_scroll)
		}
		return true
	})
	player.ScrollQueue.Clear()

	for i := range GameState.ScrollsToAdd {
		scroll_to_add := &GameState.ScrollsToAdd[i]

		server_side_scroll := ServerSideScrolls{}
		server_side_scroll.Action = scroll_to_add.Action
		server_side_scroll.Conditions = scroll_to_add.Conditions
		server_side_scroll.Posistion = GameState.OtherPlayer.POS
		server_side_scroll.ID = len(GameState.Scrolls) + 1
		server_side_scroll.Vel = scroll_to_add.Velocity

		GameState.Scrolls = append(GameState.Scrolls, server_side_scroll)
	}
	GameState.ScrollsToAdd = nil
}

func DrawScrolls(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}

	for i := range GameState.Scrolls {
		current_scroll := &GameState.Scrolls[i]
		op.GeoM.Reset()
		op.GeoM.Translate(current_scroll.Posistion.X, current_scroll.Posistion.Y)

		scroll_texture.Draw(screen, &op)
	}
}

func UpdateScrolls(Level *level.LevelStruct, Player *player.PlayerStruct) {
	for i := range GameState.Scrolls {
		if i >= len(GameState.Scrolls) {
			break
		}
		current_scroll := &GameState.Scrolls[i]

		ScrollPhysics(current_scroll, Level)

		if !IsHost {
			continue
		}

		current_scroll.Lifetime += 1

		if current_scroll.Lifetime > 100 {
			ScrollEffects(current_scroll, Player)
		}
	}

	if !IsHost {
		return
	}

	for i := len(GameState.Scrolls); i > 0; i-- {
		current_scroll := &GameState.Scrolls[i-1]
		if current_scroll.Trigged {
			utils.RemoveArrayElement(i-1, &GameState.Scrolls)
		}
	}
}

func ScrollPhysics(current_scroll *ServerSideScrolls, Level *level.LevelStruct) {
	current_scroll.Vel.Y += 0.1

	if current_scroll.Vel.X > 0.1 {
		current_scroll.Vel.X -= 0.01
	} else if current_scroll.Vel.X < -0.1 {
		current_scroll.Vel.X += 0.01
	}

	if math.Abs(current_scroll.Vel.X) < 0.1 {
		current_scroll.Vel.X = 0
	}

	if Level.Collide(utils.Vec2{X: current_scroll.Posistion.X, Y: current_scroll.Posistion.Y + current_scroll.Vel.Y}, utils.Vec2{X: 14, Y: 5}) {
		current_scroll.Vel.Y = 0
		if current_scroll.Vel.X > 0.1 {
			current_scroll.Vel.X -= 0.1
		} else if current_scroll.Vel.X < -0.1 {
			current_scroll.Vel.X += 0.1
		}
	}
	if Level.Collide(utils.Vec2{X: current_scroll.Posistion.X + current_scroll.Vel.X, Y: current_scroll.Posistion.Y}, utils.Vec2{X: 14, Y: 5}) {
		current_scroll.Vel.X = 0
	}
	current_scroll.Posistion.X += current_scroll.Vel.X
	current_scroll.Posistion.Y += current_scroll.Vel.Y
}

func ScrollEffects(current_scroll *ServerSideScrolls, Player *player.PlayerStruct) {
	all_conditions_met := true

	for _, condition := range current_scroll.Conditions {
		if !scroll.Conditions[condition](Player.Pos, GameState.OtherPlayer.POS, current_scroll.Posistion) {
			all_conditions_met = false
		}
	}

	if all_conditions_met {
		for _, action := range current_scroll.Action {
			scroll.Actions[action](Player.Pos, GameState.OtherPlayer.POS, current_scroll.Posistion)
			current_scroll.Trigged = true
		}
	}
}
