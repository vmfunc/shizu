package server

import (
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func HandleServerConn(nConn net.Conn) {
	privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	privateBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key")
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			log.Printf("Login attempt: user=%s, pass=%s\n", c.User(), string(pass))
			return nil, ssh.ErrNoAuth
		},
	}

	config.AddHostKey(private)

	_, chans, reqs, err := ssh.NewServerConn(nConn, config)

	if err != nil {
		log.Printf("Failed to establish SSH connection: %s\n", err)
		return
	}

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		newChannel.Reject(ssh.Prohibited, "no channels allowed")
	}
}
