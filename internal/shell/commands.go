package shell

import (
	"io"
)

func Whoami(i *Interpreter) bool {
	io.WriteString(i.Session, i.Session.User())
	io.WriteString(i.Session, "\n")
	return true
}
