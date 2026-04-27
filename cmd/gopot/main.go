package main

import (
	"fmt"
	"gopot/internal/config"
	"gopot/internal/db"
	"gopot/internal/shell"
	"io"
	"log"
	"net"
	"sync"
	"time"

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
	attempts := newAttemptCounter()

	server := &ssh.Server{
		Addr:        addr,
		IdleTimeout: 2 * time.Minute,
		MaxTimeout:  10 * time.Minute,
		Handler: func(s ssh.Session) {
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
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			host, _, _ := net.SplitHostPort(ctx.RemoteAddr().String())
			n := attempts.bump(host)
			ok := n > cfg.Auth.MinAttempts
			if err := session_store.InsertAuthAttempt(ctx.User(), password, host, ok); err != nil {
				log.Println("insert auth attempt:", err)
			}
			log.Printf("auth attempt #%d user=%q from=%s success=%v", n, ctx.User(), host, ok)
			return ok
		},
	}

	if cfg.Server.HostKey != "" {
		if err := server.SetOption(ssh.HostKeyFile(cfg.Server.HostKey)); err != nil {
			log.Fatalf("load host key %q: %v", cfg.Server.HostKey, err)
		}
	}

	log.Println("Starting SSH honeypot on", port, "...")
	log.Fatal(server.ListenAndServe())
}

type attemptCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func newAttemptCounter() *attemptCounter {
	return &attemptCounter{m: make(map[string]int)}
}

func (a *attemptCounter) bump(ip string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.m[ip]++
	return a.m[ip]
}
