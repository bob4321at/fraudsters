package server

import (
	"encoding/json"
	"log"
	"main/level"
	"main/player"
	"main/scroll"
	"main/utils"
	"math"
	"net"
	"time"

	"github.com/bob4321at/textures"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var Connected = false

type OtherPlayerStruct struct {
	POS    utils.Vec2
	Health int
}

type GameStateStruct struct {
	OtherPlayer  OtherPlayerStruct
	Scrolls      []ServerSideScrolls
	ScrollsToAdd []scroll.ScrollInventoryStruct
	Spells       []scroll.NetworkedSpell
}

var other_player_texture = textures.NewTexture("./art/other_player.png", "")
var scroll_texture = textures.NewTexture("./art/scroll.png", "")

var GameState = GameStateStruct{
	OtherPlayerStruct{},
	[]ServerSideScrolls{},
	[]scroll.ScrollInventoryStruct{},
	[]scroll.NetworkedSpell{},
}

var waiting_for_other = true

var Connection *net.UDPConn

var IsHost = false

var ServerAddress = "localhost:8080"

const BUFFERSIZE = 1024 * 8

var OtherPlayerDrawnPos utils.Vec2

func StartServer(Player *player.PlayerStruct) {
	log.Println("starting server")

	address, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		panic(err)
	}

	Connection, err := net.ListenUDP("udp", address)
	if err != nil {
		panic(err)
	}
	defer Connection.Close()

	buffer := make([]byte, BUFFERSIZE)

	_, clientAddr, err := Connection.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NewClient: ", string(utils.DecodeBinary(buffer)))

	IsHost = true

	first_update := false

	for {
		if waiting_for_other {
			game_state_bytes := utils.DecodeBinary(buffer)
			json.Unmarshal(game_state_bytes, &GameState)

			game_state_this_user := GameState
			game_state_this_user.OtherPlayer = OtherPlayerStruct{Player.Pos, Player.Health}

			game_state_bytes, err := json.Marshal(game_state_this_user)
			if err != nil {
				panic(err)
			}

			message := utils.EncodeBinary(game_state_bytes)

			_, err = Connection.WriteToUDP(message, clientAddr)
			if err != nil {
				panic(err)
			}
			if clientAddr != nil {
				waiting_for_other = false
				Connected = true
			}

			time.Sleep(time.Second)

			player.ScrollQueue.Clear()
		} else {
			// Get Other Player
			time.Sleep(time.Second / 15)

			buffer := make([]byte, BUFFERSIZE)

			_, clientAddr, err := Connection.ReadFromUDP(buffer)
			if err != nil {
				log.Fatal(err)
			}

			OriginalScrolls := []ServerSideScrolls{}
			for i := range GameState.Scrolls {
				copy := GameState.Scrolls[i]
				OriginalScrolls = append(OriginalScrolls, copy)
			}

			OriginalSpells := []scroll.NetworkedSpell{}
			for i := range GameState.Spells {
				if i >= len(GameState.Spells) {
					break
				}
				copy := GameState.Spells[i]
				OriginalSpells = append(OriginalSpells, copy)
			}
			other_player_real_health := GameState.OtherPlayer.Health

			other_game_state := utils.DecodeBinary(buffer)
			json.Unmarshal(other_game_state, &GameState)

			GameState.Scrolls = OriginalScrolls
			GameState.Spells = OriginalSpells
			GameState.OtherPlayer.Health = other_player_real_health

			AddScrolls(Player)
			AddSpellsToServer()

			if !first_update {
				OtherPlayerDrawnPos = GameState.OtherPlayer.POS
				first_update = true
			}

			// Send Data
			game_state_this_user := GameState
			game_state_this_user.OtherPlayer = OtherPlayerStruct{Player.Pos, GameState.OtherPlayer.Health}

			game_state_bytes, err := json.Marshal(game_state_this_user)
			if err != nil {
				panic(err)
			}

			message := utils.EncodeBinary(game_state_bytes)

			_, err = Connection.WriteToUDP(message, clientAddr)
			if err != nil {
				panic(err)
			}
		}
	}
}

func Draw(screen *ebiten.Image) {
	if Connected {
		ebitenutil.DebugPrint(screen, "Connected to server")

		op := ebiten.DrawImageOptions{}
		op.GeoM.Translate(OtherPlayerDrawnPos.X, OtherPlayerDrawnPos.Y)
		other_player_texture.Draw(screen, &op)

		DrawScrolls(screen)
		DrawSpells(screen)
	}
}

func Update(Level *level.LevelStruct, Player *player.PlayerStruct) {
	if Connected {
		UpdateScrolls(Level, Player)
		UpdateSpells(Player)

		if utils.GetDist(OtherPlayerDrawnPos, GameState.OtherPlayer.POS) > 1 {
			angle := math.Atan2(OtherPlayerDrawnPos.Y-GameState.OtherPlayer.POS.Y, OtherPlayerDrawnPos.X-GameState.OtherPlayer.POS.X)

			OtherPlayerDrawnPos.X -= math.Cos(angle) * 3 * (utils.GetDist(OtherPlayerDrawnPos, GameState.OtherPlayer.POS) / 20)
			OtherPlayerDrawnPos.Y -= math.Sin(angle) * 3 * (utils.GetDist(OtherPlayerDrawnPos, GameState.OtherPlayer.POS) / 20)

		} else {
			OtherPlayerDrawnPos = GameState.OtherPlayer.POS
		}
	}
}

func ConnectToServer(Player *player.PlayerStruct) {
	log.Println("connecting to server")

	address, err := net.ResolveUDPAddr("udp", ServerAddress)
	if err != nil {
		panic(err)
	}

	Connection, err := net.DialUDP("udp", nil, address)
	if err != nil {
		panic(err)
	}
	defer Connection.Close()

	game_state_to_send := GameState
	game_state_to_send.OtherPlayer = OtherPlayerStruct{Player.Pos, Player.Health}

	game_state_bytes, err := json.Marshal(game_state_to_send)
	if err != nil {
		panic(err)
	}

	message := utils.EncodeBinary(game_state_bytes)

	_, err = Connection.Write(message)
	if err != nil {
		panic(err)
	}

	Connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	buffer := make([]byte, BUFFERSIZE)
	_, _, err = Connection.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("Receive error: %v", err)
		return
	}

	other_game_state_bytes := utils.DecodeBinary(buffer)
	json.Unmarshal(other_game_state_bytes, &GameState)

	Connected = true
	waiting_for_other = false

	time.Sleep(time.Second)

	for {
		// SEND SCROLLS
		time.Sleep(time.Second / 15)
		player.ScrollQueue.Range(func(key, value any) bool {
			scrolls_to_add := value.(scroll.ScrollInventoryStruct)
			GameState.ScrollsToAdd = append(GameState.ScrollsToAdd, scrolls_to_add)
			return true
		})
		player.ScrollQueue.Clear()

		// SEND SELF
		game_state_to_send := GameState
		game_state_to_send.OtherPlayer = OtherPlayerStruct{Player.Pos, Player.Health}

		game_state_bytes, err := json.Marshal(game_state_to_send)
		if err != nil {
			panic(err)
		}

		message := utils.EncodeBinary(game_state_bytes)

		_, err = Connection.Write(message)
		if err != nil {
			panic(err)
		}

		// GET OTHER PLAYER
		Connection.SetReadDeadline(time.Now().Add(5 * time.Second))
		buffer := make([]byte, BUFFERSIZE)
		_, _, err = Connection.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Receive error: %v", err)
			return
		}

		other_game_state_bytes := utils.DecodeBinary(buffer)
		json.Unmarshal(other_game_state_bytes, &GameState)
		Player.Health = GameState.OtherPlayer.Health
	}
}
