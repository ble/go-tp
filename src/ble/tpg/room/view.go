package room

import (
	"ble/tpg/model"
	"ble/web"
	"encoding/json"
	"time"
)

type playerJson struct {
	model.Player
	IsYou bool
}

func (p playerJson) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["pseudonym"] = p.Pseudonym()
	result["id"] = p.Pid()
	if p.IsYou {
		result["isYou"] = true
	}
	return json.Marshal(result)
}

type drawingJsonSimple struct {
	model.Drawing
	web.Switchboard
}

func (d drawingJsonSimple) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["id"] = d.Did()
	url := d.URLOf(d.Drawing)
	if url != nil {
		result["url"] = url.String()
	} else {
		result["url"] = nil
	}
	return json.Marshal(result)
}

type stackJsonSimple struct {
	model.Stack
	web.Switchboard
}

func (s stackJsonSimple) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["id"] = s.Sid()
	url := s.URLOf(s.Stack)
	if url != nil {
		result["url"] = url.String()
	} else {
		result["url"] = nil
	}
	return json.Marshal(result)
}

type stackJson struct {
	model.Stack
	web.Switchboard
}

func (s stackJson) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["id"] = s.Sid()
	url := s.URLOf(s.Stack)
	if url != nil {
		result["url"] = url.String()
	} else {
		result["url"] = nil
	}
	drawings := s.AllDrawings()
	cDrawings := make([]drawingJsonSimple, len(drawings), len(drawings))
	for i := range drawings {
		cDrawings[i] = drawingJsonSimple{drawings[i], s.Switchboard}
	}
	result["drawings"] = cDrawings
	return json.Marshal(result)
}

type gameJson struct {
	model.Game
	lastTime time.Time
	web.Switchboard
	requestingPlayerId string
}

func (g gameJson) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["id"] = g.Gid()
	url := g.URLOf(g.Game)
	if url != nil {
		result["url"] = url.String()
	} else {
		result["url"] = nil
	}
	stacksInPlay := g.StacksInProgress()
	stacksInPlayIds := make(map[string][]string)
	for player, stacks := range stacksInPlay {
		stackIds := make([]string, len(stacks), len(stacks))
		for i, stack := range stacks {
			stackIds[i] = stack.Sid()
		}
		stacksInPlayIds[player.Pid()] = stackIds
	}
	result["stacksInPlay"] = stacksInPlayIds

	stacks := g.Stacks()
	cStacks := make([]stackJson, len(stacks), len(stacks))
	for i, stack := range stacks {
		cStacks[i] = stackJson{stack, g.Switchboard}
	}
	result["stacks"] = cStacks

	players := g.Players()
	cPlayers := make([]playerJson, len(players), len(players))
	for i, player := range players {
		cPlayers[i] = playerJson{player, g.requestingPlayerId == player.Pid()}
	}
	result["players"] = cPlayers
	result["lastTime"] = g.lastTime.UnixNano() / 1000
	result["isStarted"] = g.IsStarted()
	result["isComplete"] = g.IsComplete()
	return json.Marshal(result)
}
