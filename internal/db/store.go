package db

import (
	"database/sql"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type Store struct {
	db *sql.DB
}

// Store constructor
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertSession(se ssh.Session) (int64, error) {
	host, portStr, _ := net.SplitHostPort(se.RemoteAddr().String())
	port, _ := strconv.Atoi(portStr)

	var pubkey sql.NullString
	pk := se.PublicKey()
	if pk != nil {
		pubkey.String = strings.TrimSpace(string(gossh.MarshalAuthorizedKey(pk)))
		pubkey.Valid = true
	}

	res, err := s.db.Exec(
		`INSERT INTO sessions (user, remote_ip, remote_port, pubkey, client, started_at) VALUES (?, ?, ?, ?, ?, ?)`,
		se.User(), host, port, pubkey, se.Context().ClientVersion(), time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) EndSession(id int64) error {
	_, err := s.db.Exec(
		`UPDATE sessions SET ended_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339Nano), id,
	)
	return err
}

func (s *Store) InsertCommand(sessionID int64, command string) error {
	_, err := s.db.Exec(
		`INSERT INTO commands (session_id, command, executed_at) VALUES (?, ?, ?)`,
		sessionID, command, time.Now().UTC().Format(time.RFC3339Nano),
	)
	return err
}
