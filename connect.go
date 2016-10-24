package sshautomate

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"strconv"
)

// KeyboardInteractiveChallengePassword
// Returns an interactive challenge solver that can complete a password challenge without user input
func KeyboardInteractiveChallengePassword(password string) func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
	return func(
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
