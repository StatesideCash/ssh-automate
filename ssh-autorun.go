package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
)

// Much of this shamelessly ripped from the GoDocs because they are great

// Gets the number of users on the server
// TODO Convert return to int of the lines counted
func GetUserCount(server string, port int, config *ssh.ClientConfig) (string, error) {
	session, err := EstablishSession(server, port, config)
	if err != nil {
		return "", err
	}
	defer session.Close()

	command := "/usr/bin/who -u"
	output, err := ExecuteCommand(session, command)
	return output, err
}

// Executes a command on the session and returns the output
func ExecuteCommand(session *ssh.Session, command string) (string, error) {
	var b bytes.Buffer
	session.Stdout = &b
	err := session.Run(command)
	return b.String(), err
}

func EstablishSession(server string, port int, config *ssh.ClientConfig) (*ssh.Session, error) {
	// Establish connection to the server
	client, err := ssh.Dial("tcp", server+":"+strconv.Itoa(port), config)
	if err != nil {
		return nil, err
	}

	// Create a session for interacting with the server
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, err
}

func main() {
	log.SetFlags(log.Lshortfile | log.Llongfile)

	// Config options for how to connect
	username := "bsc4155"
	password := "BscTjc1995"
	server := "glados.cs.rit.edu"
	port := 22

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil

		}
		return []string{password}, nil

	}

	// The SSH config to connect to the server
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(keyboardInteractiveChallenge),
		},
	}

	output, err := GetUserCount(server, port, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
