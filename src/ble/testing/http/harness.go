package http

import (
	"github.com/ble/cookiejar"
	. "net/http"
	"net/http/httptest"
	"net/url"
	. "testing"
)

type ResponsePredicate func(r *Response) (bool, []interface{})
type ClientFactory func() *Client
type ServerFactory func() *httptest.Server

type SimpleCase struct {
	string
	*Request
	ResponsePredicate
}

type TestHarness struct {
	*T
	*httptest.Server
	*url.URL
}

func NewHarness(t *T, s ServerFactory) *TestHarness {
	th := new(TestHarness)
	th.T = t
	th.Server = s()
	th.URL, _ = url.Parse(th.Server.URL)
	return th
}

func (th *TestHarness) Stop() {
	th.Server.Close()
}

func (th *TestHarness) processReport(desc string, isError bool, args []interface{}) {
	if isError {
		th.T.Error(append([]interface{}{desc}, args...)...)
	} else if len(args) > 0 {
		th.T.Log(append([]interface{}{desc}, args...)...)
	}
}

func (th *TestHarness) SimpleTest(s SimpleCase, c *Client) {
	response, err := c.Do(s.Request)
	if err != nil {
		th.Fatal(err)
	} else {
		isError, args := s.ResponsePredicate(response)
		th.processReport(s.string, isError, args)
	}
}

func asSlice(args ...interface{}) []interface{} {
	return args
}

var ShouldStatusBe = func(statusCode int, equal bool) ResponsePredicate {
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

var StatusShouldBe = func(statusCode int) ResponsePredicate {
	return ShouldStatusBe(statusCode, true)
}

var StatusShouldNotBe = func(statusCode int) ResponsePredicate {
	return ShouldStatusBe(statusCode, false)
}

var CookieClient ClientFactory = func() *Client {
	client := new(Client)
	client.Jar = cookiejar.NewJar(true)
	return client
}

var PlainClient ClientFactory = func() *Client {
	client := new(Client)
	return client
}

func FromHandler(h Handler) ServerFactory {
	return ServerFactory(func() *httptest.Server {
		return httptest.NewServer(h)
	})
}
