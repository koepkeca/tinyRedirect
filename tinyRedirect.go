package tinyRedirect

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	//We use 1080 as the default port.
	DEFAULT_INPUT_PORT = "1080"
	//We are using 302 as the default redirect type
	DEFAULT_REDIRECT_TYPE = http.StatusFound //this is 302..
	//Default Redirect To, Just the programs github page
	DEFAULT_REDIRECT_TO = "https://github.com/koepkeca/tinyRedirect"
	//Environment variable for specifying listening port
	ENV_VAR_LISTEN_PORT = "LISTEN_PORT"
	//Environment variable for specifying listening address
	ENV_VAR_LISTEN_ADDR = "LISTEN_ADDR"
	//Environment variable for specifying redirect type [301-309]
	ENV_VAR_REDIR_TYPE = "REDIR_TYPE"
	//Environment variable for specifying redirect destination URL
	ENV_VAR_REDIR_DEST = "REDIR_DEST"
)

//EnvData holds the raw data from the environment
type EnvData struct {
	ListenAddr string
	PortNbr    string
	RedirType  string
	RedirTo    string
}

//NewEnvData reads the input from the environment
func NewEnvData() (e EnvData) {
	e.PortNbr = os.Getenv(ENV_VAR_LISTEN_PORT)
	e.ListenAddr = os.Getenv(ENV_VAR_LISTEN_ADDR)
	e.RedirType = os.Getenv(ENV_VAR_REDIR_TYPE)
	e.RedirTo = os.Getenv(ENV_VAR_REDIR_DEST)
	return
}

//Parse extracts a ServerConfig from a EnvData
func (e EnvData) Parse() (c ServerConfig) {
	c.ListenString = validateListenString(e.ListenAddr, e.PortNbr)
	c.RedirType = validateRedirect(e.RedirType)
	c.RedirTo = validateRedirectTarget(e.RedirTo)
	return
}

// ServerConfig contains server configuration data
type ServerConfig struct {
    ListenString string
    RedirType    int
    RedirTo      string
}

// Run is the main entry point. If you're using daemontools or 
// docker you want to make sure this method continues running
// and is not forked or run in a go routine.
func (c ServerConfig) Run() {
	http.HandleFunc("/",c.redirector)
	http.HandleFunc("/stat", c.statListener)
	err := http.ListenAndServe(c.ListenString,nil)
	log.Print("tinyRedirect starting up..")
	if err != nil {
		log.Fatal("Tinyredirect: ", err)
	}
	return
}

// Redirect sends the redirect
func (c ServerConfig) redirector(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Redirecting to: %s\n", c.RedirTo)
	http.Redirect(w, r, c.RedirTo, c.RedirType)
}

// Poll responds with a HTTP 200, useful to see if the service is up
func (c ServerConfig) statListener(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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


