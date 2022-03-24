package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Global variable to hold connection ID
var conn *websocket.Conn

// Global variable to hold connection ID
var conId *websocket.Conn

var err error

// Carries messages from server
type Message struct {
	//Id        *websocket.Conn `json:"id"` 
	Username  string          `json:"username"`
	Message   string          `json:"message"`
	//Reciptent string          `json:"reciptent"`
}

func main () {
	// Default connection data
	schema := "ws://"
	host := "localhost"
	port := ":8989"
	apiEndpoint := "/webSocket"
	Origin := GetOutboundIP() //Find machines public ip address
	server := schema + host + port + apiEndpoint

	// Get users input from args and check if correct amount of args are given and if the're suitable
	if len(os.Args) == 2 || len(os.Args) >3 {
		log.Fatal("Error: Incorrect number of arguments given. [USAGE]: nc $server-adress $port  [DEFAULT USAGE] nc")
	} else if len(os.Args) == 3{
		host = os.Args[1]
		// Check if port given in args is in range of useable ports
		if portChk(os.Args[2]) {
			port = os.Args[2]
		}
	}

	// Make the the websocket connection
	conn, err = websocket.Dial(server, "", Origin)
	if err != nil {
		log.Fatal("Error: Couldn't connecting to the chat server. HINT: Are you using the correct port? Is it switched on?")
	}

	runClient()
}

func runClient() {
	
	welcomeMsg()

	// Start a Go routine to listen for messages from server
	go getMsg()

	name := getName()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		check(scanner.Err())

		letSend := Message{Username: name, Message: scanner.Text()}
		if err := websocket.JSON.Send(conn, letSend); err != nil { //send message
			panic(err)
		}

		cmd := exec.Command("clear") //clear console then history is should arrive from server to refresh
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	defer conn.Close()
}

// Constantly scans web socket for incomming messages and prints
func getMsg() {
	text := Message{}
	for {
		if err := websocket.JSON.Receive(conn, &text); err != nil {
			fmt.Println("Error: The server has disconnected you. Hint: Is it turned on!")
			log.Fatal()
		} else {
			//conId = text.Id
			stamped := timeStamp(text)
			fmt.Println(stamped.Message)
		}
	}
}

// Takes users name and sends an initiation msg so the server can match websocket id with username
func getName() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("[ENTER YOUR NAME]: ")
	scanner.Scan()
	name := scanner.Text()
	// send init handshake msg
	letSend := Message{Username: name}
	if err := websocket.JSON.Send(conn, letSend); err != nil {
		panic(err)
	}
	return name
}

// Gets machines public ip address
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	check(err)
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	return "http://" + localAddr
}

// Prints welcome message from txt file
func welcomeMsg() {
	fileIO, err := os.OpenFile("welcome.txt", os.O_RDWR, 0600)
	check(err)
	defer fileIO.Close()
	rawBytes, err := ioutil.ReadAll(fileIO)
	check(err)
	lines := strings.Split(string(rawBytes), "\n")
	for _, i := range lines {
		fmt.Println(i)
	}
}

// Checks errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func portChk(s string) bool {
	i, _ := strconv.Atoi(s)
	if i < 1024 || i > 65352 {
		log.Fatal("Your port number is not in the correct range of 1024-65352")
	}
	return true
}

// This checks if msg username is server, if not it adds timestamp and returns the edited message
func timeStamp(msg Message) Message {
	if msg.Username != "Server" {
		currentTime := time.Now()
		time := currentTime.Format("2006.01.02 15:04:05")                    // format time
		msg.Message = "[" + time + "][" + msg.Username + "]: " + msg.Message // add timestamp to message
	}
	return msg
}
