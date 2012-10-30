package handler

import (
	"ble/game"
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
	th := NewHarness(t)
	th.start(hJoinF)
	defer th.stop()
	url := th.URL.String()

	tests := []simpleCase{}

	//All methods except for post are not allowed.
	s := StatusMethodNotAllowed
	for _, method := range []string{"POST", "GET", "HEAD", "PUT", "OPTIONS"} {
		r, _ := NewRequest(method, url, nil)
		t := simpleCase{"Allowed methods", r, ShouldStatusBe(s, method != "POST")}
		tests = append(tests, t)
	}

	//Posted body must be below 1kb
	s = StatusRequestEntityTooLarge
	for _, v := range []int{1, 128, 1023, 1024, 4096} {
		body := strings.NewReader(strings.Repeat(" ", v))
		r, _ := NewRequest("POST", url, body)
		t := simpleCase{"Post length", r, ShouldStatusBe(s, v >= 1024)}
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
			t := simpleCase{"Content type", r, ShouldStatusBe(s, v != json)}
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
		for _, v := range [][]byte{b1, b2, b3} {
			err := json.Unmarshal(v, make(map[string]interface{}))
			r, _ := NewRequest("POST", url, bytes.NewReader(v))
			r.Header["Content-Type"] = []string{"application/json"}
			t := simpleCase{"(In)valid JSON", r, ShouldStatusBe(s, err != nil)}
			tests = append(tests, t)
		}
	}

	client := plainClient()
	for _, test := range tests {
		th.simpleTest(test.string, client, test.Request, test.responsePred)
	}
}

/*
func TestJoinMethods(t *T) {
  th := NewHarness(t)
  th.start(hJoinF)
  defer th.stop()
  client := plainClient()

  for _, method := range []string{"GET", "HEAD", "PUT", "OPTIONS"} {
    request, _ := NewRequest(method, th.URL.String(), nil)
    th.simpleTest(client, request, StatusShouldBe(StatusMethodNotAllowed))
  }

  request, _ := NewRequest("POST", th.URL.String(), nil)
  th.simpleTest(client, request, StatusShouldNotBe(StatusMethodNotAllowed))
}

func TestJoinLimit(t *T) {
  th := NewHarness(t)
  th.start(hJoinF)
  defer th.stop()
  client := plainClient()

  contentLengths := []int{1, 128, 1023,  1024, 4096}
  tooLarge := StatusRequestEntityTooLarge
  for _, v := range contentLengths {
    body := strings.NewReader(strings.Repeat(" ", v))
    request, _ := NewRequest("POST", th.URL.String(), body)
    th.simpleTest(client, request, ShouldStatusBe(tooLarge, v >= 1024))
  }

}
*/
/*
func TestJoinJSON(t *T) {
  testFun := func(contentType string) harnessTest {
    return func(t *T, u url.URL, g game.GameAgent, c clientF) {
      request, _ := NewRequest("POST", u.String(), nil)
      client := c()
      response, _ := client.Do(request)
      defer response.Body.Close()
      if (len(bodyContent) < 1024) ==
         (response.StatusCode == StatusRequestEntityTooLarge) {
        t.Error(
          "Bad response code ",
          response.StatusCode,
          " for content of length ",
          len(bodyContent))
      }
    }
  }
  contentLengths := []int{1, 128, 1023,  1024, 4096}
  for _, v := range contentLengths {
    bodyContent := strings.Repeat(" ", v)
    gameTestHarness(t, hJoinF, plainClient, testFun(bodyContent))
  }

}
*/
