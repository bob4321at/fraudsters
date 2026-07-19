package server

import (
	"main/player"
	"main/scroll"
	"main/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

func AddSpellsToServer() {
	scroll.SpellQueue.Range(func(key, value any) bool {
		current_spell_to_add_data, ok := scroll.SpellQueue.Load(key)

		if ok {
			current_spell_to_add := current_spell_to_add_data.(scroll.NetworkedSpell)

			GameState.Spells = append(GameState.Spells, current_spell_to_add)
		}

		return true
	})
	scroll.SpellQueue.Clear()
}

func DrawSpells(screen *ebiten.Image) {
	for i := range GameState.Spells {
		if i < len(GameState.Spells) {
			spell := &GameState.Spells[i]
			spell.Draw(screen)
		}
	}
}

func UpdateSpells(Player *player.PlayerStruct) {
	for i := range GameState.Spells {
		if i >= len(GameState.Spells) {
			break
		}

		spell := &GameState.Spells[i]
		var info scroll.SpellInfo
		if IsHost {
			info := scroll.SpellInfo{
				HostPosistion:   Player.Pos,
				ClientPosistion: GameState.OtherPlayer.POS,
			}
			spell.Update(&info)
		} else {
			info := scroll.SpellInfo{
				HostPosistion:   GameState.OtherPlayer.POS,
				ClientPosistion: Player.Pos,
			}
			spell.Update(&info)
		}

		if !IsHost {
			continue
		}

		spell.Lifetime -= 1

		if spell.Lifetime <= 0 {
			spell.Destroy = true
		}

		if utils.Collide(Player.Pos, utils.Vec2{X: 12, Y: 16}, spell.Position, spell.Size) {
			spell.Destroy = true
			Player.Health -= spell.Damage
		}

		if utils.Collide(info.ClientPosistion, utils.Vec2{X: 12, Y: 16}, spell.Position, spell.Size) {
			spell.Destroy = true
			length := utils.GetSyncLength(&scroll.ConditionsToApply)
			scroll.ConditionsToApply.Store(length+1, scroll.Condition{Name: "DamageClient", Value: spell.Damage})
		}
	}

	if !IsHost {
		return
	}

	scroll.ConditionsToApply.Range(func(key, value any) bool {
		condition, ok := value.(scroll.Condition)
		if !ok {
			return false
		}
		Conditions[condition.Name](condition.Value)
		return true
	})
	scroll.ConditionsToApply.Clear()

	for i := len(GameState.Spells); i > 0; i-- {
		current_spell := &GameState.Spells[i-1]
		if current_spell.Destroy {
			utils.RemoveArrayElement(i-1, &GameState.Spells)
		}
	}
}
