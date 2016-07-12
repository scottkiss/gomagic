package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"github.com/scottkiss/gomagic/webmagic"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MIDDLEWARE_NAME = "session-cookie"
	SESSION_ID_KEY  = "_SID"
	TIMESTAMP_KEY   = "_STS"
	COOKIE_PREFIX   = "webmagic"
)

var (
	sessionConfig        *SessionConfig
	cookieKeyValueParser = regexp.MustCompile("\x00([^:]*):([^\x00]*)\x00")
	storedCookie         *http.Cookie
)

func init() {
	sessionConfig = new(SessionConfig)
}

type SessionConfig struct {
	CookiePrefix        string
	CookieDomain        string
	CookieSecure        bool
	ExpireAfterDuration time.Duration
	SecretKey           string
}

type Session struct {
	webmagic.Middleware
	//session data
	values map[string]string
	rw     http.ResponseWriter
	lock   sync.RWMutex
}

func (session *Session) Id() string {
	if sessionId, ok := session.values[SESSION_ID_KEY]; ok {
		return sessionId
	}
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}
	session.lock.Lock()
	defer session.lock.Unlock()
	session.values[SESSION_ID_KEY] = hex.EncodeToString(buffer)
	return session.values[SESSION_ID_KEY]
}

//return cookie containing the signed session.
func (session *Session) Cookie() *http.Cookie {
	var sessionValue string
	t := session.getExpirationTime()
	session.lock.Lock()
	defer session.lock.Unlock()
	if session.values == nil {
		session.values = make(map[string]string)
	}
	session.values[TIMESTAMP_KEY] = getSessionExpirationCookie(t)
	for k, v := range session.values {
		if strings.ContainsAny(k, ":\x00") {
			panic("Session keys may not have colons or null bytes")
		}
		if strings.Contains(v, "\x00") {
			panic("Session values may not have null bytes")
		}
		sessionValue += "\x00" + k + ":" + v + "\x00"
	}
	sessionData := url.QueryEscape(sessionValue)
	if sessionConfig.CookiePrefix == "" {
		sessionConfig.CookiePrefix = COOKIE_PREFIX
	}
	return &http.Cookie{
		Name:     sessionConfig.CookiePrefix + "_SESSION",
		Value:    sign(sessionData, sessionConfig.SecretKey) + "-" + sessionData,
		Domain:   sessionConfig.CookieDomain,
		Path:     "/",
		HttpOnly: true,
		Secure:   sessionConfig.CookieSecure,
		Expires:  t.UTC(),
	}
}

func (session *Session) SetResponseWriter(rw http.ResponseWriter) {
	session.rw = rw
}

func (session *Session) Set(key, value string) {
	restoredSession := getSessionFromCookie(session.Cookie(), sessionConfig.SecretKey)
	restoredSession.values[key] = value
	cookie := restoredSession.Cookie()
	if session.rw == nil {
		panic("session rw property is nil,must call SetResponseWriter(rw http.ResponseWriter) method")
	}
	storedCookie = cookie
	http.SetCookie(session.rw, cookie)
}

func (session *Session) Get(key string) string {
	restoredSession := getSessionFromCookie(storedCookie, sessionConfig.SecretKey)
	return restoredSession.values[key]
}

func (session *Session) restoreSession(req *http.Request, config *SessionConfig) *Session {
	cookie, err := req.Cookie(config.CookiePrefix + "_SESSION")
	if err != nil {
		if session.values == nil {
			session.values = make(map[string]string)
		}
		session.Id()
		sessionConfig.CookiePrefix = config.CookiePrefix
		sessionConfig.CookieDomain = config.CookieDomain
		sessionConfig.CookieSecure = config.CookieSecure
		sessionConfig.ExpireAfterDuration = config.ExpireAfterDuration
		sessionConfig.SecretKey = config.SecretKey
		return session
	} else {
		storedCookie = cookie
		return getSessionFromCookie(cookie, sessionConfig.SecretKey)
	}
}

func getSessionFromCookie(cookie *http.Cookie, secretKey string) *Session {
	session := new(Session)
	session.values = make(map[string]string)
	hyphen := strings.Index(cookie.Value, "-")
	if hyphen == -1 || hyphen >= len(cookie.Value)-1 {
		return session
	}

	sig, data := cookie.Value[:hyphen], cookie.Value[hyphen+1:]
	//verify the signature
	if !verifySign(data, sig, secretKey) {
		log.Println("verify cookie signature failed")
		return session
	}
	//parse session key-value
	val, _ := url.QueryUnescape(data)
	//log.Println(val)
	if matchs := cookieKeyValueParser.FindAllStringSubmatch(val, -1); matchs != nil {
		for _, match := range matchs {
			session.values[match[1]] = match[2]
		}
	}
	if sessionTimeoutExpiredOrMissing(session) {
		session = new(Session)
	}
	return session
}

func sessionTimeoutExpiredOrMissing(session *Session) bool {
	if exp, currentSession := session.values[TIMESTAMP_KEY]; !currentSession {
		return true
	} else if exp == "session" {
		return false
	} else if expInt, _ := strconv.Atoi(exp); int64(expInt) < time.Now().Unix() {
		return true
	}
	return false
}

func (session *Session) getExpirationTime() time.Time {
	if sessionConfig.ExpireAfterDuration == 0 || session.values[TIMESTAMP_KEY] == "session" {
		return time.Time{}
	}
	return time.Now().Add(sessionConfig.ExpireAfterDuration)
}

func sign(message, secretKey string) string {
	mac := hmac.New(sha1.New, []byte(secretKey))
	io.WriteString(mac, message)
	return hex.EncodeToString(mac.Sum(nil))
}

func verifySign(message, sig, secretKey string) bool {
	return hmac.Equal([]byte(sig), []byte(sign(message, secretKey)))
}

func getSessionExpirationCookie(t time.Time) string {
	if t.IsZero() {
		return "session"
	}
	return strconv.FormatInt(t.Unix(), 10)
}

//init session middleware
func Init(config *SessionConfig) *webmagic.Middleware {
	s := new(webmagic.Middleware)
	s.Name = MIDDLEWARE_NAME
	s.Handler = func(ctx *webmagic.Context, middleware *webmagic.Middleware) {
		session := new(Session)
		restoredSession := session.restoreSession(ctx.Request, config)
		http.SetCookie(ctx.ResponseWriter, restoredSession.Cookie())
		middleware.Next(ctx)
	}
	return s
}
