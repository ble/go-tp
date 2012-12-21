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

func asJoinGame(data []byte) (*action, error) {
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

func JoinGame(playerId, name string) action {
	return action{
		ActionType: "joinGame",
		Who:        playerId,
		Name:       name}
}

func Chat(playerId, content string) action {
	return action{
		ActionType: "chat",
		Who:        playerId,
		Content:    content}
}

func PassStack(playerId, recipientId, stackId string) action {
	return action{
		ActionType: "passStack",
		Who:        playerId,
		ToWhom:     recipientId,
		StackId:    stackId}
}
