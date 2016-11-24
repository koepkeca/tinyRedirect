package main

import (
	"github.com/koepkeca/tinyRedirect"
)

func main() {
	env := tinyRedirect.NewEnvData()
	sc := env.Parse()
	sc.Run()
	return
}
