package server

var Conditions = map[string]func(value int){
	"DamageClient": DamageClient,
}

func DamageClient(value int) {
	GameState.OtherPlayer.Health -= value
}
