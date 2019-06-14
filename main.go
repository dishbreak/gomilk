package main

import (
	"fmt"

	"github.com/dishbreak/gomilk/api/auth"
)

func main() {
	frob, err := auth.GetFrob()
	if err != nil {
		panic(err)
	}
	fmt.Println(frob)
}
