package utilmagic

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	HTTPMETHOD = map[string]string{
		"DELETE": "DELETE",
		"HEAD":   "HEAD",
		"GET":    "GET",
		"POST":   "POST",
		"PUT":    "PUT",
	}
)

const (
	HEADER_CONTENT_ENCODING = "Content-Encoding"

	VERSION           = "0.0.1"
	DEFAULT_USERAGENT = "gomagic httpclient - " + VERSION

	PROXY_HTTP    = "HTTP"
	PROXY_SOCKS4  = "SOCKS4"
	PROXY_SOCKS5  = "SOCKS5"
	PROXY_SOCKS4A = "SOCKS4A"

	DEFAULT_ST_CONNECTTIMEOUT = 60
	DEFAULT_ST_TIMEOUT        = 60
	DEFAULT_ST_COOKIEJAR      = true
)

var defaultSetting = map[string]interface{}{
	"ST_USERAGENT":      DEFAULT_USERAGENT,
	"ST_COOKIEJAR":      DEFAULT_ST_COOKIEJAR,
	"ST_CONNECTTIMEOUT": DEFAULT_ST_CONNECTTIMEOUT,
	"ST_TIMEOUT":        DEFAULT_ST_TIMEOUT,
}

var defaultHeaders = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36",
	"Accept-Encoding": "gzip, deflate, sdch",
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
	"Accept-Language": "zh-CN,zh;q=0.8,en;q=0.6,zh-TW;q=0.4",
}

type Resp struct {
	*http.Response
}

type Settings map[string]interface{}
type Headers map[string]string

func (r *Resp) ReadBytes() ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch r.Header.Get(HEADER_CONTENT_ENCODING) {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		if err != nil {
			return nil, err
		}
	default:
		reader = r.Body
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func (r *Resp) ReadString() (string, error) {
	bytes, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func DefaultHttpClinet() *HttpClient {
	return CustomHttpClient(defaultSetting, defaultHeaders)

}

func CustomHttpClient(settings Settings, headers Headers) *HttpClient {
	client := &HttpClient{
		Headers:        headers,
		Settings:       settings,
		reuseTransport: true,
		reuseCookieJar: true,
	}
	return client
}

type HttpClient struct {
	Headers        map[string]string
	Settings       map[string]interface{}
	Transport      http.RoundTripper
	reuseTransport bool
	reuseCookieJar bool
	CookieJar      http.CookieJar
	cookies        []*http.Cookie
	lock           *sync.Mutex
}

func (h *HttpClient) reset() {
	h.reuseTransport = true
	h.reuseCookieJar = true
	h.cookies = nil
	if h.lock != nil {
		h.lock.Unlock()
	}
}

func (h *HttpClient) Request(method string, url string, headers map[string]string,
	body io.Reader) (*Resp, error) {
	var err error
	var transport http.RoundTripper
	var jar http.CookieJar
	if _, ok := headers["User-Agent"]; !ok {
		if headers == nil {
			headers = make(map[string]string)

		}
		headers["User-Agent"] = defaultSetting["ST_USERAGENT"].(string)
	}

	headers = MapMerge(h.Headers, headers)
	if h.Transport == nil || !h.reuseTransport {
		transport, err = initTransport(h.Settings)
		if err != nil {
			h.reset()
			return nil, err
		}
		if h.reuseTransport {
			h.Transport = transport
		}
	} else {
		transport = h.Transport
	}

	if h.CookieJar == nil || !h.reuseCookieJar {
		jar, err = initJar(h.Settings)
		if err != nil {
			h.reset()
			return nil, err
		}

		if h.reuseCookieJar {
			h.CookieJar = jar
		}
	} else {
		jar = h.CookieJar
	}

	h.reset()

	c := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	req, err := newRequest(method, url, headers, body)
	if err != nil {
		return nil, err
	}

	if jar != nil {
		jar.SetCookies(req.URL, h.cookies)
	} else {
		for _, cookie := range h.cookies {
			req.AddCookie(cookie)
		}
	}
	res, err := c.Do(req)
	return &Resp{res}, err
}

func (h *HttpClient) Head(url string, params map[string]string) (*Resp,
	error) {
	url = concatParams(url, params)
	return h.Request(HTTPMETHOD["HEAD"], url, nil, nil)
}

func (h *HttpClient) Get(url string, params map[string]string) (*Resp,
	error) {
	url = concatParams(url, params)
	return h.Request(HTTPMETHOD["GET"], url, nil, nil)
}

func (h *HttpClient) Delete(url string, params map[string]string) (*Resp,
	error) {
	url = concatParams(url, params)
	return h.Request(HTTPMETHOD["DELETE"], url, nil, nil)
}

func (h *HttpClient) Put(url string, params map[string]string) (*Resp, error) {
	url = concatParams(url, params)
	return h.Request(HTTPMETHOD["PUT"], url, nil, nil)
}

func (h *HttpClient) Post(url string, params map[string]string) (*Resp,
	error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	body := strings.NewReader(MapToString(params))
	return h.Request(HTTPMETHOD["POST"], url, headers, body)
}

func newRequest(method string, url string, headers map[string]string,
	body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

func initTransport(settings map[string]interface{}) (http.RoundTripper, error) {
	transport := &http.Transport{}
	connectTimeout := 0
	if connectTimeout_, ok := settings["ST_CONNECTTIMEOUT"]; ok {
		if connectTimeout, ok := connectTimeout_.(int); ok {
			connectTimeout = connectTimeout * 1000
		} else {
			return nil, fmt.Errorf("ST_CONNECTTIMEOUT must be int")
		}
	}
	timeout := 0
	if timeout_, ok := settings["ST_TIMEOUT"]; ok {
		if timeout, ok := timeout_.(int); ok {
			timeout = timeout * 1000
		} else {
			return nil, fmt.Errorf("ST_TIMEOUT must be int")
		}
	}
	if timeout > 0 && (connectTimeout > timeout || connectTimeout == 0) {
		connectTimeout = timeout
	}
	transport.Dial = func(network, addr string) (net.Conn, error) {
		var conn net.Conn
		var err error
		if connectTimeout > 0 {
			conn, err = net.DialTimeout(network, addr, time.Duration(connectTimeout)*time.Millisecond)
			if err != nil {
				return nil, err
			}
		} else {
			conn, err = net.Dial(network, addr)
			if err != nil {
				return nil, err
			}
		}
		if timeout > 0 {
			conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
		}
		return conn, nil
	}

	if proxyFunc_, ok := settings["ST_PROXY_FUNC"]; ok {
		if proxyFunc, ok := proxyFunc_.(func(*http.Request) (int, string, error)); ok {
			transport.Proxy = func(req *http.Request) (*url.URL, error) {
				proxyType, u_, err := proxyFunc(req)
				if err != nil {
					return nil, err
				}
				if proxyType != PROXY_HTTP {
					return nil, fmt.Errorf("only PROXY_HTTP is currently supported")
				}
				u_ = "http://" + u_
				u, err := url.Parse(u_)
				if err != nil {
					return nil, err
				}
				return u, nil
			}
		} else {
			return nil, fmt.Errorf("ST_PROXY_FUNC is not a desired function")
		}
	} else {
		var proxytype int
		if proxytype_, ok := settings["ST_PROXYTYPE"]; ok {
			if proxytype, ok = proxytype_.(int); !ok || proxytype != PROXY_HTTP {
				return nil, fmt.Errorf("ST_PROXYTYPE must be int, and only PROXY_HTTP is currently supported")
			}
		}
		var proxy string
		if proxy_, ok := settings["ST_PROXY_ADDR"]; ok {
			if proxy, ok = proxy_.(string); !ok {
				return nil, fmt.Errorf("ST_PROXY_ADDR must be string")
			}
			proxy = "http://" + proxy
			proxyUrl, err := url.Parse(proxy)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	return transport, nil
}

func initJar(settings map[string]interface{}) (http.CookieJar, error) {
	var jar http.CookieJar
	var err error
	if cookieJar, ok := settings["ST_COOKIEJAR"]; ok {
		if cookieJarBool, ok := cookieJar.(bool); ok {
			if cookieJarBool {
				jar, err = cookiejar.New(nil)
				if err != nil {
					return nil, err
				}
			}
		} else if cookieJarObj, ok := cookieJar.(http.CookieJar); ok {
			jar = cookieJarObj
		} else {
			return nil, fmt.Errorf("invalid cookiejar")
		}
	}

	return jar, nil
}

func concatParams(url string, params map[string]string) string {
	if len(params) == 0 {
		return url
	}

	if !strings.Contains(url, "?") {
		url += "?"
	}

	if strings.HasSuffix(url, "?") || strings.HasSuffix(url, "&") {
		url += MapToString(params)
	} else {
		url += "&" + MapToString(params)
	}

	return url
}
