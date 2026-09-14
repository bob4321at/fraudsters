package scroll

import (
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

func PlayerTwoHundredClose(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool {
	client_dist := utils.GetDist(ScrollPos, OtherPlayerPos)
	host_dist := utils.GetDist(ScrollPos, PlayerPos)

	return (client_dist < 200 || host_dist < 200)
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

func BlueSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/blue.json", ScrollPos)
}

func RedSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/red.json", ScrollPos)
}

func PurpleSpell(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/purple.json", ScrollPos)
}

func SmallHeal(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) {
	AddSpell("./assets/spells/heal.json", ScrollPos)
}

var Conditions = map[string]func(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2) bool{
	"start_active":             StartActive,
	"player_two_hundred_close": PlayerTwoHundredClose,
	"player_hundred_close":     PlayerHundredClose,
	"player_fifty_close":       PlayerFiftyClose,
}

var Actions = map[string]func(PlayerPos, OtherPlayerPos, ScrollPos utils.Vec2){
	"blue":       BlueSpell,
	"red":        RedSpell,
	"purple":     PurpleSpell,
	"small_heal": SmallHeal,
}
