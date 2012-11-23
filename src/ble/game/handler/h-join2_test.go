package handler

import (
	"ble/game"
	. "ble/testing/http"
	"bytes"
	"encoding/json"
	. "net/http"
	"strings"
	. "testing"
)

//"Request tests"
//JSON must parse

//Request and response tests
//JSON must be of correct format
//Name must be acceptable 
//Name must not already be in use
//Artist ID cookie must not already be set to something in use

//Response tests
//Cookie must be set after good post

func hJoinF(g game.GameAgent) Handler { return handlerJoin{g} }

func TestJoinSimple(t *T) {
	g := game.NewGame()
	th := NewHarness(t, FromHandler(handlerJoin{g}))
	defer th.Stop()
	url := th.URL.String()

	tests := []SimpleCase{}

	//All methods except for post are not allowed.
	s := StatusMethodNotAllowed
	for _, method := range []string{"POST", "GET", "HEAD", "PUT", "OPTIONS"} {
		r, _ := NewRequest(method, url, nil)
		t := SimpleCase{"Allowed methods", r, ShouldStatusBe(s, method != "POST")}
		tests = append(tests, t)
	}

	//Posted body must be below 1kb
	s = StatusRequestEntityTooLarge
	for _, v := range []int{1, 128, 1023, 1024, 4096} {
		body := strings.NewReader(strings.Repeat(" ", v))
		r, _ := NewRequest("POST", url, body)
		t := SimpleCase{"Post length", r, ShouldStatusBe(s, v >= 1024)}
		tests = append(tests, t)
	}

	//Posted content-type must be JSON
	s = StatusUnsupportedMediaType
	{
		json := "application/json"
		form := "application/x-www-form-urlencoded"
		multipart := "multipart/form-data"
		for _, v := range []string{json, form, multipart} {
			r, _ := NewRequest("POST", url, nil)
			r.Header["Content-Type"] = []string{v}
			t := SimpleCase{"Content type", r, ShouldStatusBe(s, v != json)}
			tests = append(tests, t)
		}
	}

	//JSON must parse
	s = StatusBadRequest
	{
		j1 := []interface{}{1, 2, 3, "abc"}
		j2 := map[string]interface{}{"greeting": "hello"}
		b1, _ := json.Marshal(j1)
		b2, _ := json.Marshal(j2)
		b3 := []byte("{\"kinda\":\"json\",\"not\":really}")
		for _, vv := range [][]byte{b1, b2, b3} {
			v := vv
			var z *interface{} = new(interface{})
			err := json.Unmarshal(v, z)
			t.Log(*z, err, err != nil)
			r, _ := NewRequest("POST", url, bytes.NewReader(v))
			r.Header["Content-Type"] = []string{"application/json"}
			t := SimpleCase{"(In)valid JSON", r, ShouldStatusBe(s, err != nil)}
			tests = append(tests, t)
		}
	}

	client := PlainClient()
	for _, test := range tests {
		th.SimpleTest(test, client)
	}
}
