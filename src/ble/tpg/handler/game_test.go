package handler

import (
	"ble/testing/http"
	"ble/tpg/persistence"
	"ble/tpg/room"
	"ble/tpg/switchboard"
	"bytes"
	"encoding/json"
	. "net/http"
	"os"
	"runtime/debug"
	. "testing"
)

func dieOnErr(e error, t *T) {
	if e != nil {
		debug.PrintStack()
		t.Fatal(e)
	}
}

var testDbFileName string = "test-db"

func TestGameHandler(t *T) {
	//set up working directory
	//so that we can fetch things from our static asset directory
	cwd, err := os.Getwd()
	dieOnErr(err, t)
	dieOnErr(os.Chdir(cwd+"/../../../../"), t)
	defer os.Chdir(cwd)

	//set up app backend
	backend, err := persistence.NewBackend(testDbFileName)
	dieOnErr(err, t)
	defer os.Remove(testDbFileName)

	//create handler-related stuff
	rooms := room.NewRoomService(switchboard.NewSwitchboard(), backend)
	gh := StripPrefix("/game/", &gameHandler{rooms})

	//create domain objects
	user0, _ := backend.CreateUser("a", "sd", "f")
	user0, _ = backend.LogInUser("sd", "f")
	user1, _ := backend.CreateUser("f", "ds", "a")
	user1, _ = backend.LogInUser("ds", "a")

	game, _ := backend.CreateGame("fofoyang")

	//create test server
	harness := http.NewHarness(t, http.FromHandler(gh))
	defer harness.Stop()

	//fake logging in users on their respective clients
	client0 := http.CookieClient()
	client0Cookie := &Cookie{
		Name:     "userId",
		Value:    user0.Uid(),
		Path:     "/",
		HttpOnly: true}
	client0.Jar.SetCookies(harness.URL, []*Cookie{client0Cookie})

	client1 := http.CookieClient()
	client1Cookie := &Cookie{
		Name:     "userId",
		Value:    user1.Uid(),
		Path:     "/",
		HttpOnly: true}
	client1.Jar.SetCookies(harness.URL, []*Cookie{client1Cookie})

	join0Json := `{"actionType":"joinGame","name":"dazzler"}`
	join0, err := NewRequest(
		"POST",
		harness.URL.String()+"/game/"+game.Gid()+"/join",
		bytes.NewReader([]byte(join0Json)))
	for _, cookie := range client0.Jar.Cookies(join0.URL) {
		join0.AddCookie(cookie)
	}
	respJoin0, err := client0.Do(join0)
	dieOnErr(err, t)
	j0, _ := json.Marshal(respJoin0)
	t.Log("Joining game", string(j0))
	/*
	  //This way and the above work; client.Do() does not automatically add cookies
	  respJoin0, err := client0.Post(
	    harness.URL.String()+"/"+game.Gid()+"/join",
	    "application/json",
	    bytes.NewReader([]byte(join0Json)))
	*/
	redirectHeaderValues := respJoin0.Header["Location"]
	if redirectHeaderValues == nil || len(redirectHeaderValues) != 1 {
		t.Fatal("bad Location header", redirectHeaderValues)
	}
	clientUrl, _ := respJoin0.Request.URL.Parse(redirectHeaderValues[0])
	t.Log(clientUrl.String())

	respGetBeforeJoin, err := client1.Get(clientUrl.String())
	jgbj, _ := json.Marshal(respGetBeforeJoin)
	t.Log("Attempt to access client prior to joining game: ", string(jgbj))

	join1Json := `{"actionType":"joinGame","name":"wig-fuckin' fairy folk"}`
	respJoin1, err := client1.Post(
		harness.URL.String()+"/game/"+game.Gid()+"/join",
		"application/json",
		bytes.NewReader([]byte(join1Json)))
	dieOnErr(err, t)
	j1, _ := json.Marshal(respJoin1)
	t.Log("Second player joining game", string(j1))

	respGetClient, err := client1.Get(clientUrl.String())
	j2, _ := json.Marshal(respGetClient)
	t.Log("Second player gets game client", string(j2))
}
