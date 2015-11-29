package webmagic

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Input struct {
	r *http.Request
}

func NewInput(r *http.Request) *Input {
	return &Input{r: r}
}

func (in *Input) ReadJson(v interface{}) error {
	body, err := ioutil.ReadAll(in.r.Body)
	in.r.Body.Close()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (in *Input) ReadXml(v interface{}) error {
	body, err := ioutil.ReadAll(in.r.Body)
	in.r.Body.Close()
	if err != nil {
		return err
	}
	return xml.Unmarshal(body, v)
}
