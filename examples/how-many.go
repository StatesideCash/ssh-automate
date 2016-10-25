package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/statesidecash/sshautomate"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
)

var (
	server          string
	username        string
	password        string
	manual_password bool
	port            int
)

func main() {
	// Configure logging output to be more verbose
	// This program is still in beta after all
	log.SetFlags(log.Lshortfile | log.Llongfile)

	// Config options for how to connect
	// TODO Add shorthand variants
	// TODO Add key auth
	// TODO Put the flag initialization in a separate function to declutter main
	// TODO Support for blank passwords
	flag.StringVar(&server, "server", "", "Server to connect to")
	flag.StringVar(&username, "user", "", "User account on the server")
	flag.StringVar(&password, "pass", "", "Password to authenticate with")
	flag.BoolVar(&manual_password, "manual-password", false, "Manually enter password to via prompt")
	flag.IntVar(&port, "port", 22, "Port the SSH daemon is running on")

	flag.Parse()

	// Satisfy some of the flags initializations

	if manual_password {
		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
		password = string(pass)
	}

	// Sanity checking
	if username == "" {
		log.Fatalln("No username given. Please specify --user")
	} else if server == "" {
		log.Fatalln("No server given. Please specify --server")
	} else if password == "" {
		log.Fatalln("No password given. Please specify --pass or --manual-password")
	}

	// The SSH config to connect to the server
	// TODO Maybe do something instead of throwing every auth method we have at it?
	//      Some kind of prioritization maybe? /shrug
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(sshautomate.KeyboardInteractiveChallengePassword(password)),
		},
	}

	// Establishes a session to the server
	session, err := sshautomate.EstablishSession(server, port, config)
	if err != nil {
		log.Fatal(err)
	}

	// Executes a command on the session and saves the output
	output, err := sshautomate.ExecuteCommand(session, "/usr/bin/who -u")
	if err != nil {
		log.Fatal(err)
	}

	// Counts the number of lines in the output of who -u
	counter := 0
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		counter++
	}
	fmt.Printf("Number of users: %d\n", counter)
}
