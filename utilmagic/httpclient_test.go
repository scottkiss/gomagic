package utilmagic

import (
	"testing"
)

func TestGet(t *testing.T) {
	var client *HttpClient = DefaultHttpClinet()
	res, err := client.Get("https://github.com/", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res.ReadString())
}
