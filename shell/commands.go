package shell

import (
	"io"
)

func Whoami(i *Interpreter) bool {
	io.WriteString(i.Session, "user")
	io.WriteString(i.Session,"\n")
	return true
}
