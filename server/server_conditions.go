package server

import (
	"main/player"
	"main/utils"
)

var Conditions = map[string]func(value int){
	"DamageClient":   DamageClient,
	"LaunchPlayerUp": ApplyForceToClient,
}

var PlayerRef = &player.PlayerStruct{}

func ApplyForceToClient(value int) {
	if value == 0 {
		PlayerRef.Vel.Y = -5
	}
	if value == 1 {
		GameState.VelocityToAddToClient = utils.Vec2{X: 0, Y: -5}
		GameState.ClientVelocityForcesAdded = false
		panic(GameState.VelocityToAddToClient)
	}
}

func DamageClient(value int) {
	GameState.OtherPlayer.Health -= value
}
