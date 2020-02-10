package log

/*
	Development Debgging Goodness
*/

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func Dump(args ...interface{}) {

	spew.Config.DisableMethods = true
	spew.Config.Indent = "  "

	buffer := spew.Sdump(args...)

	fmt.Println(buffer)
}
