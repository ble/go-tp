package handler

import (
	"ble/testing/http"
	"ble/tpg/persistence"
	"ble/tpg/switchboard"
	"os"
	. "testing"
)

func TestEphCreateUser(t *T) {
	backend, _ := persistence.NewBackend(testDbFileName)
	defer os.Remove(testDbFileName)
	sb := switchboard.NewSwitchboard(backend)
	harness := http.NewHarness(t, http.FromHandler(sb))
	client0 := http.CookieClient()

	accessURL := sb.URLOf(
		sb.GetEphemera().NewCreateUser(
			"binjermon",
			"benjaminster@gmail.com",
			"whackadoodle"))
	t.Log(accessURL.String())
	baseURL := harness.URL
	absoluteURL := baseURL.ResolveReference(accessURL)
	resp, err := client0.Get(absoluteURL.String())
	t.Log(resp, err)
	cJar := client0.Jar
	cookies := cJar.Cookies(baseURL)
	if len(cookies) == 0 {
		t.Fatal("no cookies set")
	}
}
