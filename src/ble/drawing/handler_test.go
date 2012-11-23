package drawing

import (
	. "ble/testing/http"
	"bytes"
	"io/ioutil"
	. "net/http"
	. "testing"
)

var logBodyAndPass ResponsePredicate = func(r *Response) (bool, []interface{}) {
	bytes, _ := ioutil.ReadAll(r.Body)
	return false, []interface{}{string(bytes)}
}

func TestDrawHandler(t *T) {
	th := NewHarness(t, FromHandler(drawingHandler{NewDrawingHandle()}))
	defer th.Stop()
	url := th.URL.String()
	client := PlainClient()

	reqPostJSON, _ := NewRequest(
		"POST",
		url,
		bytes.NewReader([]byte(testDrawParts[0])))
	testPost := SimpleCase{"Post draw part", reqPostJSON, StatusShouldBe(StatusOK)}

	reqGetJSON, _ := NewRequest("GET", url, nil)
	testGet := SimpleCase{"Get drawing", reqGetJSON, logBodyAndPass}
	th.SimpleTest(testPost, client)
	th.SimpleTest(testGet, client)
}
