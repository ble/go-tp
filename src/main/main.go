package main

import (
	"ble/game"
	"ble/game/handler"
	"log"
	. "net/http"
)

func main() {
	handlerGame := handler.NewRoomHandler(game.NewGame())
	log.Fatal(ListenAndServe(":24769", handlerGame))
}
