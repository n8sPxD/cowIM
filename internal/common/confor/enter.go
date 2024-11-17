package confor

import (
	"fmt"
	"os"
)

func FromEnv(conf any, env string) {
	appEnv := os.Getenv(env)
	if appEnv == "" {
		panic(fmt.Sprintf("environment variable %s not set", env))
	}
	conf = appEnv
}