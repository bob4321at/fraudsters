package server

import "fmt"

var Conditions = map[string]func(value int){
	"DamageClient": DamageClient,
}

func DamageClient(value int) {
	GameState.OtherPlayer.Health -= value
	fmt.Println("testing")
	fmt.Println(GameState.OtherPlayer.Health)
}
