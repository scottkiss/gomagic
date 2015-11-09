package httpmagic

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"path/filepath"
	"strconv"
)

const (
	applicationJSON = "application/json"
	applicationXML  = "application/xml"
	textXML         = "text/xml"
	textHTML        = "text/html"
)

type Output struct {
	w http.ResponseWriter
	r *http.Request
}

func NewOutput(w http.ResponseWriter, r *http.Request) *Output {
	return &Output{w: w, r: r}
}

func (out *Output) Header(key, val string) {
	out.r.Header.Set(key, val)
}

func (out *Output) Body(b []byte) {
	out.Header("Content-Length", strconv.Itoa(len(b)))
	out.w.Write(b)
}

//Output json data
func (out *Output) Json(v interface{}, hasIndent bool) {
	var (
		content []byte
		err     error
	)

	if hasIndent {
		content, err = json.MarshalIndent(v, "", "  ")
	} else {
		content, err = json.Marshal(v)
	}

	if err != nil {
		http.Error(out.w, err.Error(), http.StatusInternalServerError)
		return
	}
	out.Header("Content-Length", strconv.Itoa(len(content)))
	out.Header("Content-Type", applicationJSON)
	out.w.Write(content)
}

//Output xml data
func (out *Output) Xml(v interface{}) {
	content, err := xml.Marshal(v)
	if err != nil {
		http.Error(out.w, err.Error(), http.StatusInternalServerError)
		return
	}
	out.Header("Content-Length", strconv.Itoa(len(content)))
	out.Header("Content-Type", textXML+"; charset=utf-8")
	out.w.Write(content)
}

//Output file
func (out *Output) File(file string) {
	out.Header("Content-Description", "File Transfer")
	out.Header("Content-Type", "application/octet-stream")
	out.Header("Content-Disposition", "attachment; filename="+filepath.Base(file))
	out.Header("Content-Transfer-Encoding", "binary")
	out.Header("Expires", "0")
	out.Header("Cache-Control", "must-revalidate")
	out.Header("Pragma", "public")
	http.ServeFile(out.w, out.r, file)
}

//Output Html data
func (out *Output) Html(v interface{}) {
	out.Header("Content-Type", textHTML+"; charset=utf-8")
	out.Body(v.([]byte))
}

func (out *Output) ServeAccept(v interface{}) {
	accept := out.r.Header.Get("Accept")
	switch accept {
	case applicationJSON:
		out.Json(v, true)
	case applicationXML, textXML:
		out.Xml(v)
	case textHTML:
		out.Html(v)
	default:
		out.Json(v, true)
	}
	return
}
