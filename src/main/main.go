package main

import (
	"ble/tpg/persistence"
	"ble/tpg/switchboard"
	"fmt"
	"net"
	. "net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {
	tmpDbFile := "db"
	port := ":24769"
	hostname := "localhost"
	base := "http://" + hostname + port

	fileAfter := fmt.Sprintf("%s-%s", tmpDbFile, time.Now())
	backend, _ := persistence.NewBackend(tmpDbFile)
	defer os.Rename(tmpDbFile, fileAfter)
	sb := switchboard.NewSwitchboard(backend)

	game, _ := backend.CreateGame("thangdangdoodle")
	gameUrlBase := sb.URLOf(game)
	gameJoinUrl, _ := url.Parse(base + gameUrlBase.Path + "join-client")

	ephemera := sb.GetEphemera()
	ephemeralJoin := ephemera.NewCreateUser("benjaminstarrr", "the.bomb@thebomb.com", "asdfqweruxl", gameJoinUrl)
	ephUrl, _ := url.Parse(base + sb.URLOf(ephemeralJoin).Path)
	fmt.Println(ephUrl.String())

	server := Server{
		Addr:         port,
		Handler:      sb,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second}
	listener, _ := net.Listen("tcp", port)
	go server.Serve(listener)
	chInterrupt := make(chan os.Signal)
	signal.Notify(chInterrupt, os.Interrupt)
	select {
	case <-chInterrupt:
		return
	}
}
