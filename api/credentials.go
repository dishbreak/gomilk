package api

import (
	"errors"
	"os"
)

var (
	APIKey       = os.Getenv("GOMILK_API_KEY")
	SharedSecret = os.Getenv("GOMILK_SHARED_SECRET")
)

func ValidateCredentials() {
	for _, varName := range []string{"GOMILK_API_KEY", "GOMILK_SHARED_SECRET"} {
		if _, found := os.LookupEnv(varName); !found {
			panic(errors.New("miaaing env var: " + varName))
		}
	}
}
