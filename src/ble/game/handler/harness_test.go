package handler

import (
	"ble/game"
	. "net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	. "testing"
)

type responsePred func(r *Response) (bool, []interface{})
type handlerF func(g game.GameAgent) Handler
type clientF func() *Client

type simpleCase struct {
	string
	*Request
	responsePred
}

type testHarness struct {
	*T
	*httptest.Server
	game.GameAgent
	*url.URL
}

func NewHarness(t *T) *testHarness {
	th := new(testHarness)
	th.T = t
	return th
}

func (th *testHarness) start(h handlerF) {
	th.GameAgent = game.NewGame()
	th.Server = httptest.NewServer(h(th.GameAgent))
	th.URL, _ = url.Parse(th.Server.URL)

}

func (th *testHarness) stop() {
	th.Server.Close()
	th.GameAgent.Shutdown()
}

func (th *testHarness) processReport(desc string, isError bool, args []interface{}) {
	if isError {
		th.T.Error(append([]interface{}{desc}, args...)...)
	} else if len(args) > 0 {
		th.T.Log(append([]interface{}{desc}, args...)...)
	}
}

func (th *testHarness) simpleTest(desc string, c *Client, r *Request, t responsePred) {
	response, err := c.Do(r)
	if err != nil {
		th.T.Fatal(err)
	} else {
		okay, args := t(response)
		th.processReport(desc, okay, args)
	}
}

func asSlice(args ...interface{}) []interface{} {
	return args
}

var ShouldStatusBe = func(statusCode int, equal bool) responsePred {
	return func(r *Response) (bool, []interface{}) {
		if r.StatusCode == statusCode && !equal {
			return true, asSlice("Status should not be ", statusCode)
		}
		if r.StatusCode != statusCode && equal {
			return true,
				asSlice(
					"Status should be ", statusCode,
					", is ", r.StatusCode)
		}
		return false, asSlice()
	}
}

var StatusShouldBe = func(statusCode int) responsePred {
	return ShouldStatusBe(statusCode, true)
}

var StatusShouldNotBe = func(statusCode int) responsePred {
	return ShouldStatusBe(statusCode, false)
}

var cookieClient clientF = func() *Client {
	client := new(Client)
	client.Jar = cookiejar.NewDefaultJar()
	return client
}

var plainClient clientF = func() *Client {
	client := new(Client)
	return client
}
