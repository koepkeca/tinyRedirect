package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	//We use 1080 as the default port.
	DEFAULT_INPUT_PORT    = "1080"
    //We are using 302 as the default redirect type
	DEFAULT_REDIRECT_TYPE = http.StatusFound //this is 302..
    //Default Redirect To, Just the programs github page
	DEFAULT_REDIRECT_TO   = "https://github.com/koepkeca/tinyRedirect"
)

//Global to hold the parsed system configuration
var SysConf ServerConfig

//EnvData holds the raw data from the environment
type EnvData struct {
	ListenAddr string
	PortNbr    string
	RedirType  string
	RedirTo    string
}

//NewEnvData reads the input from the environment
//NewEnvDataが環境変数を読む
func NewEnvData() (e EnvData) {
	e.PortNbr = os.Getenv("LISTEN_PORT")
	e.ListenAddr = os.Getenv("LISTEN_ADDR")
	e.RedirType = os.Getenv("REDIR_TYPE")
	e.RedirTo = os.Getenv("REDIR_DEST")
	return
}

//Parse extracts a ServerConfig from a EnvData
func (e EnvData) Parse() (c ServerConfig) {
	c.ListenString = validateListenString(e.ListenAddr, e.PortNbr)
	c.RedirType = validateRedirect(e.RedirType)
	c.RedirTo = validateRedirectTarget(e.RedirTo)
	return
}

// validateListenString validates and creates the connect string
func validateListenString(a string, p string) (as string) {
	if a != "" {
		as = a
	}
	as += ":"
	_, err := strconv.Atoi(p)
	if err != nil || p == "" {
		fmt.Printf("Port number missing or invalid, defaulting to %s\n", DEFAULT_INPUT_PORT)
		as += DEFAULT_INPUT_PORT
		return
	}
	as += p
	return
}

//validateRedirect validates the redirect type
func validateRedirect(t string) (rc int) {
	rc = DEFAULT_REDIRECT_TYPE
	code, err := strconv.Atoi(t)
	if err != nil {
		fmt.Printf("Invalid redirect, defaulting to %d\n", rc)
		return
	}
	if code < 300 || code > 309 {
		fmt.Printf("Redirect code out of bounds, defaulting to %d\n", rc)
		return
	}
	rc = code
	return
}

//validateRedirectTarget validates the redirected target URL
func validateRedirectTarget(t string) (r string) {
	r = DEFAULT_REDIRECT_TO
	if t != "" {
		r = t
	}
	return
}

// ServerConfig contains server configuration data
type ServerConfig struct {
	ListenString string
	RedirType    int
	RedirTo      string
}

func Poll(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Redirecting to: %s\n", SysConf.RedirTo)
	http.Redirect(w, r, SysConf.RedirTo, SysConf.RedirType)
}

func main() {
	env := NewEnvData()
	SysConf = env.Parse()
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/poll", Poll)
	err := http.ListenAndServe(SysConf.ListenString, nil)
	fmt.Printf("Starting up\n")
	if err != nil {
		log.Fatal("tinyredirect:", err)
	}
}
