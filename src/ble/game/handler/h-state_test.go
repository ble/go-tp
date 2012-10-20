package handler

import (
	"ble/game"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	. "testing"
)

func Test_State_Handler(t *T) {
	agent := game.NewGame()
	defer agent.Shutdown()

	a0, _ := agent.AddArtist("Sammy")
	agent.AddArtist("Jimbo")
	agent.Start()

	server := httptest.NewServer(handlerState{agent})
	defer server.Close()

	submitter := http.Client{}

	request, _ := http.NewRequest("GET", server.URL, nil)
	response, _ := submitter.Do(request)
	respBodyBytes, _ := ioutil.ReadAll(response.Body)
	respBody := string(respBodyBytes)
	t.Log(respBody)

	aIdCookie := &http.Cookie{
		Name:  "artistId",
		Value: a0.Id}
	request.AddCookie(aIdCookie)

	response, _ = submitter.Do(request)
	respBodyBytes, _ = ioutil.ReadAll(response.Body)
	respBody = string(respBodyBytes)
	t.Log(respBody)
	t.Log(response)

}
