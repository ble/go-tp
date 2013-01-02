package handler

import (
	"ble/testing/http"
	"ble/tpg/persistence"
	"ble/tpg/switchboard"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	. "net/http"
	"net/url"
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
	sb := switchboard.NewSwitchboard(backend)

	//create test server
	harness := http.NewHarness(t, http.FromHandler(sb))
	defer harness.Stop()
	game, _ := backend.CreateGame("fofoyang")

	client0 := http.CookieClient()
	client0.CheckRedirect = func(req *Request, via []*Request) error {
		return errors.New("no redirects followed")
	}
	//sadpath: get join client before user is logged-in
	joinClientUrl := harness.URL.String() +
		"/game/" +
		game.Gid() +
		"/join-client"
	{
		resp, _ := client0.Get(joinClientUrl)
		r0, _ := json.Marshal(resp)
		t.Log("Sadpath: get client for joining game before user logs on", string(r0))
	}

	//happypath: create user and log in using ephemeral URL
	loginDestination, _ := url.Parse("/")
	createUser0URL := sb.URLOf(
		sb.GetEphemera().NewCreateUser(
			"binjermon",
			"benjaminster@gmail.com",
			"whackadoodle",
			loginDestination))
	createUser0URL = harness.URL.ResolveReference(createUser0URL)
	respLogin, err := client0.Get(createUser0URL.String())
	j, _ := json.Marshal(respLogin)
	t.Log("Login response:", err, string(j))

	createUser1URL := sb.URLOf(
		sb.GetEphemera().NewCreateUser(
			"spamban",
			"the.bomb@thebomb.com",
			"rah diggah",
			loginDestination))
	createUser1URL = harness.URL.ResolveReference(createUser1URL)
	respLogin, err = client1.Get(createUser1URL.String())

	//get join client after user is logged-in
	{
		resp, _ := client0.Get(joinClientUrl)
		r0, _ := json.Marshal(resp)
		t.Log("Happypath: get client for joining game after logon", string(r0))
	}

	client1 := http.CookieClient()
	client1Cookie := &Cookie{
		Name:     "userId",
		Value:    user1.Uid(),
		Path:     "/",
		HttpOnly: true}
	client1.Jar.SetCookies(harness.URL, []*Cookie{client1Cookie})

	join0Json := `{"actionType":"joinGame","name":"dazzler"}`
	respJoin0, err := client0.Post(
		harness.URL.String()+"/game/"+game.Gid()+"/join",
		"application/json",
		bytes.NewReader([]byte(join0Json)))

	dieOnErr(err, t)
	j0, _ := json.Marshal(respJoin0)
	t.Log("Joining game", string(j0))

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

	//get game state before starting
	respGetState, err := client0.Get(harness.URL.String() + "/game/" + game.Gid())
	j2, _ = json.Marshal(respGetState)
	t.Log("First player gets game state", string(j2))
	stateBody, _ := ioutil.ReadAll(respGetState.Body)
	t.Log("Game state from response:", string(stateBody))

	//have players chat
	chatJson := `{"actionType":"chat","content":"foobaf"}`
	respChat, err := client0.Post(
		harness.URL.String()+"/game/"+game.Gid()+"/chat",
		"application/json",
		bytes.NewReader([]byte(chatJson)))
	j2, _ = json.Marshal(respChat)
	t.Log("First player chats", string(j2))

	//get game events
	respEvents, err := client0.Get(
		harness.URL.String() + "/game/" + game.Gid() + "/events")
	j2, _ = json.Marshal(respEvents)
	t.Log("First gets events", string(j2))
	eventBody, _ := ioutil.ReadAll(respEvents.Body)
	t.Log("Events body", string(eventBody))

	//start game
	respStartGame, err := client1.Post(
		harness.URL.String()+"/game/"+game.Gid()+"/start",
		"application/json",
		bytes.NewReader([]byte(`{"actionType":"startGame"}`)))
	j2, _ = json.Marshal(respStartGame)
	t.Log("Second player starts game", string(j2))

	//get game state after starting
	respGetState, err = client0.Get(harness.URL.String() + "/game/" + game.Gid())
	j2, _ = json.Marshal(respStartGame)
	t.Log("First player gets game state", string(j2))
	stateBody, _ = ioutil.ReadAll(respGetState.Body)
	t.Log("Game state from response:", string(stateBody))

	//TODO: have players pass stacks until game is over

	//get game events after starting
	respEvents, err = client0.Get(
		harness.URL.String() + "/game/" + game.Gid() + "/events")
	j2, _ = json.Marshal(respEvents)
	t.Log("First gets events", string(j2))
	eventBody, _ = ioutil.ReadAll(respEvents.Body)
	t.Log("Events body", string(eventBody))
}
