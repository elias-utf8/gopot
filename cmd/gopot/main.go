package main

import (
	"database/sql"
	"gopot/internal/config"
	"gopot/internal/db"
	"gopot/internal/shell"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	port := 2222

	cfg := config.LoadConfig()

	database, err := db.Open("gopot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, cfg.Shell.Banner)

		sessionID, err := insertSession(database, s)
		if err != nil {
			log.Println("insert session:", err)
			return
		}
		log.Println("logged new connection, session id =", sessionID)
		defer endSession(database, sessionID)

		MyShell := shell.NewInterpreter(s, cfg)
		MyShell.Run()
	})

	log.Println("starting ssh server on", port, "...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func insertSession(database *sql.DB, s ssh.Session) (int64, error) {
	host, portStr, _ := net.SplitHostPort(s.RemoteAddr().String())
	port, _ := strconv.Atoi(portStr)

	var pubkey sql.NullString
	if pk := s.PublicKey(); pk != nil {
		pubkey.String = strings.TrimSpace(string(gossh.MarshalAuthorizedKey(pk)))
		pubkey.Valid = true
	}

	res, err := database.Exec(
		`INSERT INTO sessions (user, remote_ip, remote_port, pubkey, started_at) VALUES (?, ?, ?, ?, ?)`,
		s.User(), host, port, pubkey, time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func endSession(database *sql.DB, id int64) {
	if _, err := database.Exec(`UPDATE sessions SET ended_at = ? WHERE id = ?`, time.Now().UTC().Format(time.RFC3339Nano), id); err != nil {
		log.Println("end session:", err)
	}
}
