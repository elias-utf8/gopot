package main

 import (
	"decoy/shell"
 	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"os"
	"encoding/json"
	"net"
	"golang.org/x/term"
	)

type LogConnection struct {
	User string `json:"user"`
	RemoteAddr net.Addr `json:"addr"`
	PubKey ssh.PublicKey `json:"pubkey"`
}

func createFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, err
}

func main() {
	port := 2222
	var file string = "log.json"
	
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Welcome to Ubuntu 22.04.3 LTS (GNU/Linux 5.15.0-91-generic x86_64)\n")
		
		connection := LogConnection{User: s.User(), RemoteAddr: s.RemoteAddr(), PubKey: s.PublicKey()}

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

		MyShell := shell.NewShell(s)

		// Input loop
		_, _, isPty := s.Pty()
		if isPty {
		    term := term.NewTerminal(s, "user@srv-prod-01:~$ ")
		    for {
		        line, err := term.ReadLine()
		        if err != nil {
		            break
		        }
		        if line == "whoami" {
		        	_ = MyShell.Whoami()
		        }
		        	
		    }
		}
				
	})

	log.Println("starting ssh server on", port, "...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
 }
