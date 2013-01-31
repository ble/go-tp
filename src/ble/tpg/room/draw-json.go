package room

import (
	"ble/tpg/drawing"
	"encoding/json"
	"io"
)

type drawAction struct {
	ActionType string           `json:"actionType"`
	Content    drawing.DrawPart `json:"content"`
}

func readDrawAction(r io.Reader) (*drawAction, error) {
	result := &drawAction{}
	result.Content = drawing.DefaultDrawPart
	err := json.NewDecoder(r).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func asDrawAction(data []byte) (*drawAction, error) {
	result := &drawAction{}
	result.Content = drawing.DefaultDrawPart
	err := json.Unmarshal(data, result)
	return result, err
}
