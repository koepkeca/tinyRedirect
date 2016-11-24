package tinyRedirect

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	//DEFAULT_INPUT_PORT is the default setting for the listening port
	DEFAULT_INPUT_PORT = "1080"
	//DEFAULT_REDIRECT_TYPE is the default setting for the http redirect code [301-309]
	DEFAULT_REDIRECT_TYPE = http.StatusFound //this is 302..
	//DEFAULT_REDIRECT_TO is the default URL to redirect to
	DEFAULT_REDIRECT_TO = "https://github.com/koepkeca/tinyRedirect"
	//ENV_VAR_LISTEN_PORT is the name of the environment variable containing the listening port
	ENV_VAR_LISTEN_PORT = "LISTEN_PORT"
	//ENV_VAR_LISTEN_ADDR is the name of the environment variable containing the listening address
	ENV_VAR_LISTEN_ADDR = "LISTEN_ADDR"
	//ENV_VAR_REDIR_TYPE is the name of the environment variable containing the redirect type
	ENV_VAR_REDIR_TYPE = "REDIR_TYPE"
	//ENV_VAR_REDIR_DEST is the name of the environment variable  specifying redirect destination URL
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
	c.Logger = log.New(os.Stdout, "tr:", log.Ldate|log.Ltime|log.Lmicroseconds)
	c.ListenString = c.validateListenString(e.ListenAddr, e.PortNbr)
	c.RedirType = c.validateRedirect(e.RedirType)
	c.RedirTo = c.validateRedirectTarget(e.RedirTo)
	return
}

//ServerConfig contains server configuration data
type ServerConfig struct {
	ListenString string
	RedirType    int
	RedirTo      string
	Logger       *log.Logger
}

//Run is the main entry point. If you're using daemontools or
//docker you want to make sure this method continues running
//and is not forked or run in a go routine.
func (c ServerConfig) Run() {
	http.HandleFunc("/", c.redirector)
	http.HandleFunc("/stat", c.statListener)
	err := http.ListenAndServe(c.ListenString, nil)
	c.Logger.Print("tinyRedirect starting up..")
	if err != nil {
		c.Logger.Fatal("Tinyredirect: ", err)
	}
	return
}

//Redirect sends the redirect
func (c ServerConfig) redirector(w http.ResponseWriter, r *http.Request) {
	c.Logger.Printf("Redirecting to: %s\n", c.RedirTo)
	http.Redirect(w, r, c.RedirTo, c.RedirType)
	return
}

//Poll responds with a HTTP 200, useful to see if the service is up
func (c ServerConfig) statListener(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

//validateListenString validates and creates the connect string
func (c ServerConfig) validateListenString(a string, p string) (as string) {
	if a != "" {
		as = a
	}
	as += ":"
	_, err := strconv.Atoi(p)
	if err != nil || p == "" {
		c.Logger.Printf("Port number missing or invalid, defaulting to %s\n", DEFAULT_INPUT_PORT)
		as += DEFAULT_INPUT_PORT
		return
	}
	as += p
	return
}

//validateRedirect validates the redirect type
func (c ServerConfig) validateRedirect(t string) (rc int) {
	rc = DEFAULT_REDIRECT_TYPE
	code, err := strconv.Atoi(t)
	if err != nil {
		c.Logger.Printf("Invalid redirect, defaulting to %d\n", rc)
		return
	}
	if code < 300 || code > 309 {
		c.Logger.Printf("Redirect code out of bounds, defaulting to %d\n", rc)
		return
	}
	rc = code
	return
}

//validateRedirectTarget validates the redirected target URL
func (c ServerConfig) validateRedirectTarget(t string) (r string) {
	r = DEFAULT_REDIRECT_TO
	if t != "" {
		r = t
	}
	return
}
