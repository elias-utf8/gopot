package shell

import (
	"github.com/gliderlabs/ssh"
	"io"
)

type Shell struct {
	 Session ssh.Session
}

// Shell constructor
func NewShell(s ssh.Session) *Shell {
	return &Shell{Session: s}
}

func (shell *Shell) Whoami() bool {
	io.WriteString(shell.Session, "user")
	io.WriteString(shell.Session,"\n")
	return true
}
