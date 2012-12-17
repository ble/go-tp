package handler

import ()

type Action struct {
	ActionType string `json:"actionType"`
	Name       string `json:"name,omitempty"`
}
