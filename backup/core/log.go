package core

import (
	"fmt"

	"github.com/dhindustries/graal"
)

func voidlog(api *graal.Api, f string, v ...interface{}) {
}

func logf(api *graal.Api, f string, v ...interface{}) {
	fmt.Printf(f, v...)
}
