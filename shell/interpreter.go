package shell

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type Interpreter struct {
	 Session ssh.Session
}

// Interpreter constructor
func NewInterpreter(s ssh.Session) *Interpreter {
	return &Interpreter{Session: s}
}

func (i *Interpreter) Run() {
	// Input loop
	_, _, isPty := i.Session.Pty()
	if isPty {
		term := term.NewTerminal(i.Session, "user@srv-prod-01:~$ ")
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			if line == "whoami" {
				_ = Whoami(i)
			}
		}
	}		
}
