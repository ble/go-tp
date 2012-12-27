package room

import (
	"encoding/json"
	"errors"
)

type action struct {
	ActionType string `json:"actionType"`
	Who        string `json:"who,omitempty"`
	ToWhom     string `json:"toWhom,omitempty"`
	StackId    string `json:"stackId,omitempty"`
	Content    string `json:"content,omitempty"`
	Name       string `json:"name,omitempty"`
}

func asjoingame(data []byte) (*action, error) {
	result := &action{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	if result.ActionType == "joinGame" &&
		result.Who == "" &&
		result.Name != "" {
		return result, nil
	}
	return nil, errors.New("bad input json")
}

func joingame(playerid, name string) action {
	return action{
		ActionType: "joinGame",
		Who:        playerid,
		Name:       name}
}

func aschat(data []byte) (*action, error) {
	result := &action{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	if result.ActionType == "chat" &&
		result.Content != "" {
		return result, nil
	}
	return nil, errors.New("bad input json")
}
func chat(playerid, content string) action {
	return action{
		ActionType: "chat",
		Who:        playerid,
		Content:    content}
}

func aspassstack(data []byte) (*action, error) {
	result := &action{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	if result.ActionType == "passStack" {
		return result, nil
	}
	return nil, errors.New("bad input json")
}

func passstack(playerid, recipientid, stackid string) action {
	return action{
		ActionType: "passStack",
		Who:        playerid,
		ToWhom:     recipientid,
		StackId:    stackid}
}

func asstartgame(data []byte) (*action, error) {
	result := &action{}
	err := json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	if result.ActionType == "startGame" {
		return result, nil
	}
	return nil, errors.New("bad input json")
}

func startgame(playerid string) action {
	return action{
		ActionType: "startGame",
		Who:        playerid}
}
