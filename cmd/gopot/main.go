package main

import (
	"gopot/internal/config"
	"gopot/internal/db"
	"gopot/internal/shell"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
)

func main() {
	port := 2222

	cfg := config.LoadConfig()

	database, err := db.Open("gopot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	store := db.NewStore(database)

	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, cfg.Shell.Banner)

		sessionID, err := store.InsertSession(s)
		if err != nil {

			log.Println("insert session:", err)
			return
		}
		log.Println("logged new connection, session id =", sessionID)
		defer store.EndSession(sessionID)

		MyShell := shell.NewInterpreter(s, cfg)
		MyShell.Run()
	})

	log.Println("starting ssh server on", port, "...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
