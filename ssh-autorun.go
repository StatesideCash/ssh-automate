package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
)

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
	// Configure logging output to be more verbose
	log.SetFlags(log.Lshortfile | log.Llongfile)

	// Config options for how to connect
	// TODO Add shorthand variants
	// TODO Add key auth
	// TODO Put the flag initialization in a separate function to declutter main
	// TODO Support for blank passwords
	server := flag.String("server", "", "Server to connect to")
	username := flag.String("user", "", "User account on the server")
	password := flag.String("pass", "", "Password to authenticate with")
	manual_password := flag.Bool("manual-password", false, "Manually enter password to via prompt")
	port := flag.Int("port", 22, "Port the SSH daemon is running on")
	command := flag.String("cmd", "", "The command to execute on the server")

	flag.Parse()

	// Satisfy some of the flags initializations

	if *manual_password {
		// Set the password from a prompt on STDIN (preferably with hidden output)
		//var err error
		//fmt.Print("Password: ")
		//read := bufio.NewReader(os.Stdin)
		//*password, err = read.ReadString('\n')
		//if err != nil {
		//log.Fatal(err)
		//}
		fmt.Print("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
		*password = string(pass)
	}

	// Sanity checking
	if *username == "" {
		log.Fatalln("No username given. Please specify --user")
	} else if *server == "" {
		log.Fatalln("No server given. Please specify --server")
	} else if *command == "" {
		log.Fatalln("No command given. Please specify --cmd")
	} else if *password == "" {
		log.Fatalln("No password given. Please specify --pass or --manual-password")
	}
	////////////////////////////

	// KeyboardInteractiveChallengePassword
	// Authenticates to the server by automatically providing a password
	// to a keyboard-interactive server so you don't have to type it out
	var KeyboardInteractiveChallengePassword = func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil

		}
		return []string{*password}, nil
	}

	// Actual "logic"
	// The SSH config to connect to the server
	// TODO Maybe do something instead of throwing every auth method we have at it?
	//      Some kind of prioritization maybe? /shrug
	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(*password),
			ssh.KeyboardInteractive(KeyboardInteractiveChallengePassword),
		},
	}

	// Establishes a session to the server
	session, err := EstablishSession(*server, *port, config)
	if err != nil {
		log.Fatal(err)
	}

	// Executes a command on the session and saves the output
	output, err := ExecuteCommand(session, *command)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}
