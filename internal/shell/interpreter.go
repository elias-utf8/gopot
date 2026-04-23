package shell

import (
	"gopot/internal/config"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type Interpreter struct {
	 Session ssh.Session
	 Cfg *config.Config
}

// Interpreter constructor
func NewInterpreter(s ssh.Session, c *config.Config) *Interpreter {
	return &Interpreter{Session: s, Cfg: c}
}

func (i *Interpreter) Run() {
	// Input loop
	_, _, isPty := i.Session.Pty()
	if isPty {
		term := term.NewTerminal(i.Session, i.Cfg.Shell.Prompt)
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
