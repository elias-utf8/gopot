package shell

import (
	"fmt"
	"gopot/internal/config"
	"gopot/internal/db"
	"log"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type Interpreter struct {
	Session    ssh.Session
	Session_id int64
	Cfg        *config.Config
	Store      *db.Store
}

// Interpreter constructor
func NewInterpreter(s ssh.Session, sid int64, c *config.Config, st *db.Store) *Interpreter {
	return &Interpreter{Session: s, Session_id: sid, Cfg: c, Store: st}
}

func (i *Interpreter) Run() {
	// Input loop
	_, _, isPty := i.Session.Pty()
	if isPty {
		prompt := fmt.Sprintf("%s%s", i.Session.User(), i.Cfg.Shell.Prompt)
		term := term.NewTerminal(i.Session, prompt)
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			if line == "whoami" {
				_ = Whoami(i)
				i.Store.InsertCommand(i.Session_id, "whoami")
				log.Println("logged new command, session id =", i.Session_id)

			}
		}
	}
}
