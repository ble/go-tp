package room

import (
	"ble/tpg/model"
	"ble/web"
	"encoding/json"
)

type playerJson struct{ model.Player }

func (p playerJson) MarshalJSON() ([]byte, error) {
	result := make(map[string]interface{})
	result["pseudonym"] = p.Pseudonym()
	result["id"] = p.Pid()
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
	web.Switchboard
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
	cStacksInPlay := make(map[string][]stackJsonSimple)
	for player, stacks := range stacksInPlay {
		cStacks := make([]stackJsonSimple, len(stacks), len(stacks))
		for i, stack := range stacks {
			cStacks[i] = stackJsonSimple{stack, g.Switchboard}
		}
		cStacksInPlay[player.Pid()] = cStacks
	}
	result["stacksInPlay"] = cStacksInPlay

	stacks := g.Stacks()
	cStacks := make([]stackJson, len(stacks), len(stacks))
	for i, stack := range stacks {
		cStacks[i] = stackJson{stack, g.Switchboard}
	}
	result["stacks"] = cStacks

	players := g.Players()
	cPlayers := make([]playerJson, len(players), len(players))
	for i, player := range players {
		cPlayers[i] = playerJson{player}
	}
	result["players"] = cPlayers
	return json.Marshal(result)
}
