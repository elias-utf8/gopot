package main

 import (
     "github.com/gliderlabs/ssh"
     "io"
     "log"
     "os"
     "encoding/json"
     "net"
 )

type Connection struct {
	User string `json:"user"`
	RemoteAddr net.Addr `json:"addr"`
	PubKey ssh.PublicKey `json:"pubkey"`	
}

func createFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, err
}

func main() {
	port := 2222
	var file string = "log.json"
	
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "┌────────────────────────────────────────────────────┐\n")
		io.WriteString(s, "│     UNAUTHORIZED ACCESS IS STRICTLY PROHIBITED     │\n")
		io.WriteString(s, "└────────────────────────────────────────────────────┘\n")
		
		connection := Connection{User: s.User(), RemoteAddr: s.RemoteAddr(), PubKey: s.PublicKey()}

		file, err := createFile(file)
		if err != nil {
			log.Fatal(err)
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(connection); err != nil {
			panic(err)
		}

		log.Println("logged new connection")
		defer file.Close() 
	})

	log.Println("starting ssh server on", port, "...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
 }
