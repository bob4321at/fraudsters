package scroll

import (
	"fmt"
	"main/utils"
	"sync"
)

type ScrollInventoryStruct struct {
	Conditions []string
	Action     []string
	Velocity   utils.Vec2
}

var SpellQueue sync.Map

func NewInvScroll(Conditions []string, Actions []string) ScrollInventoryStruct {
	new_scroll := ScrollInventoryStruct{}

	new_scroll.Conditions = Conditions
	new_scroll.Action = Actions

	return new_scroll
}

func StartActive(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool {
	return true
}

func PlayerHundredClose(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool {
	client_dist := utils.GetDist(ScrollPos, OtherPlayerPos)
	host_dist := utils.GetDist(ScrollPos, PlayerPos)

	return (client_dist < 100 || host_dist < 100)
}

func PlayerFiftyClose(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool {
	client_dist := utils.GetDist(ScrollPos, OtherPlayerPos)
	host_dist := utils.GetDist(ScrollPos, PlayerPos)

	return (client_dist < 50 || host_dist < 50)
}

func PlayerTenClose(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool {
	client_dist := utils.GetDist(ScrollPos, OtherPlayerPos)
	host_dist := utils.GetDist(ScrollPos, PlayerPos)

	return (client_dist < 10 || host_dist < 10)
}

func PrintSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	fmt.Println("testing spells")
}

func BlueSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/blue.json", ScrollPos)
}

func RedSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/red.json", ScrollPos)
}

func PurpleSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/blue.json", ScrollPos)
	AddSpell("./assets/spells/red.json", ScrollPos)
	AddSpell("./assets/spells/purple.json", ScrollPos)
}

var Conditions = map[string]func(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool{
	"start_active":         StartActive,
	"player_hundred_close": PlayerHundredClose,
	"player_fifty_close":   PlayerFiftyClose,
	"player_ten_close":     PlayerTenClose,
}

var Actions = map[string]func(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2){
	"print_spell": PrintSpell,
	"blue":        BlueSpell,
	"red":         RedSpell,
	"purple":      PurpleSpell,
}
