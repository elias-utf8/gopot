package main

import (
	"fmt"
	"gopot/internal/config"
	"gopot/internal/db"
	"gopot/internal/shell"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
)

func main() {

	cfg := config.LoadConfig()
	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)

	fmt.Println(cfg.Server.Banner)

	database, err := db.Open("gopot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	session_store := db.NewStore(database)

	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, cfg.Shell.Banner)

		sessionID, err := session_store.InsertSession(s)
		if err != nil {

			log.Println("insert session:", err)
			return
		}
		log.Println("logged new connection, session id =", sessionID)
		defer session_store.EndSession(sessionID)

		shell_store := db.NewStore(database)
		MyShell := shell.NewInterpreter(s, sessionID, cfg, shell_store)
		MyShell.Run()
	})

	log.Println("Starting SSH honeypot on", port, "...")
	log.Fatal(ssh.ListenAndServe(addr, nil))
}
