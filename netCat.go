package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"golang.org/x/net/websocket"
)

var clients = make(map[*websocket.Conn]bool)              // map for connected clients turns false when they leave
var clientUserName = make(map[*websocket.Conn]string)     // used to match the websocket connection data to username
var usernameConnection = make(map[string]*websocket.Conn) // used to find a msgs username and send response to there connection
var mainHistory string                                    // Holds a record of the chat for new clients to use on first connection
var clientHistory = make(map[string]string)               // holds a copy of each clients chat history
var broadcast = make(chan Message)                        // channel to handle new msgs to be sent to client devices
var historyChan = make(chan Message)                      // channel to handle msgs that are being added to a clients history record

// for reading json to
type Message struct {
	Id        *websocket.Conn `json:"id"` 
	Client    string		  `json "client"`
	Username  string          `json:"username"`
	Message   string          `json:"message"`
	Reciptent string          `json:"reciptent"`
}

func main() {

	port := ":8989" // Default port

	if len(os.Args) > 2 {
		log.Fatal("[USAGE]: ./TCPChat $port")
	} else if len(os.Args) == 2 {
		// Check if port given in args is in range of useable ports
		if portChk(os.Args[1]) {
			port = ":" + os.Args[1]
		}
	}

	// file server to send the web version of the chat on request
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// api endpoint for websocket connection, calls HandleConnections func
	http.Handle("/webSocket", websocket.Handler(handleConnections))

	// Start listening on both channels for messages, concurrently
	go handleMessages()

	go addHistory()

	// Start the server
	log.Println("http server started on " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil { // check errors
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(sockID *websocket.Conn) {
	// Make sure connection closes
	defer sockID.Close()
	// Register new client in map of clients so can loop through and send msgs to them all
	clients[sockID] = true

	for { // listen for msgs
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := websocket.JSON.Receive(sockID, &msg)
		msg.Id = sockID

		// if client has left delete them from list of clients and send message notifting users they left
		if err != nil {
			log.Printf(clientUserName[sockID] + " has left")            // print msg to server terminal
			msg.Username = "Server"                                     // create a server msg
			msg.Reciptent = "forAll"                                    // set to for all so every reciptent gets the msg
			msg.Message = clientUserName[sockID] + " has left our chat" // add the message saying user has left to a Message object
			delete(clients, sockID)                                     // delete them from list of clients so we know longer broadcast msgs to them
			broadcast <- msg                                            // send the message object to be broadcast to devices
			historyChan <- msg                                          // send to history channel to be added to users history
			break
		}

		// detect initiation msg and match username to session id
		if len(clientUserName[sockID]) < 1 { // check if username has already been matched to websocket connection data
			clientUserName[sockID] = msg.Username     // if hasnt been matched already then match it
			usernameConnection[msg.Username] = sockID // create a reveresed map so we can find the connection that matches a username.
			clientHistory[msg.Username] = ""
			fmt.Println("Handshake succeeded: matched user " + clientUserName[sockID] + " to session id")
			if len(mainHistory) > 0 && msg.Client != "webapp" { // if there is a chat history send it to new client
				sendHist := Message{Message: mainHistory}
				err = websocket.JSON.Send(sockID, sendHist) // // send history to client because there screen has been cleared by their client
				clientHistory[msg.Username] = mainHistory   // add history to map for each clients history
			}
			// Send message to all users notifying them a new user has joined the chat
			msg.Message = msg.Username + " Has joined the chat"
			msg.Username = "Server"
			msg.Reciptent = "forAll"
			broadcast <- msg
			historyChan <- msg

		// Detect empty message and only add to the users screen 
		} else if len(msg.Message) < 1 {
			sendHitstory(sockID, msg.Client)		    // Send user there history because client has cleared there screen
			msg.Reciptent = clientUserName[sockID]      // send to user who sent empty msg
			broadcast <- msg                            // send to broadcast channel
			historyChan <- msg                          // send to history channel

		// Handle messages that break the max character limit of 90 by responding to user and not sending the message to be broadcast
		} else if len(msg.Message) > 90 {
			sendHitstory(sockID, msg.Client)						// Send user there history because client has cleared there screen
			msg.Reciptent = clientUserName[sockID]      // address message to user who broke the rules
			msg.Message = "Hey " + clientUserName[sockID] + ", you broke the rules! Keep your messages limited to 90 characters!"
			msg.Username = "Server" // server responds so set username for no timestamp
			broadcast <- msg        // send
			historyChan <- msg      // send to history channel

		// Detect if user is asking the server to do something and action it
		} else if strings.Contains(msg.Message, "hey server change my name to") {
			newName := strings.SplitAfter(msg.Message, "to ")
			sendHitstory(sockID, msg.Client)						// Send user there history because client has cleared there screen
			clientUserName[sockID] = newName[1]                     // update name
			clientHistory[newName[1]] = clientHistory[msg.Username] // create new username in map and copy history
			usernameConnection[newName[1]] = sockID					// add to map for sending individual messages
			msg.Reciptent = "forAll"                                // send to all reciptents
			broadcast <- msg                                        // send to broadcast channel
			historyChan <- msg                                      // send to history channel
			msg.Username = "Server"                                 // server responds so set username for no timestamp
			msg.Message = "Server: No problem, " + clientUserName[sockID] + " I changed your name!"
			broadcast <- msg   // send to broadcast channel
			historyChan <- msg // send to history channel

		// Handle normal message where user has been matched to socketID and is already in records
		} else if len(clientUserName[sockID]) > 0 {
			sendHitstory(sockID, msg.Client)						// Send user there history because client has cleared there screen
			msg.Reciptent = "forAll"                    // send to all reciptents
			broadcast <- msg                            // send to broadcast channel
			historyChan <- msg                          // send to history channel
		}

	}
}

// Listens for messages in the broadcast channel and sends them to appropriate user devices over the websocket
func handleMessages() {
	for {
		// Listen on the channel and take new msgs
		Channelmsg := <-broadcast
		// check if the msg needs timestamp and edit it
		msg := Channelmsg

		if msg.Reciptent == "forAll" { //
			// Send it out to every client that is currently connected
			for client := range clients { // loop through all clients
				err := websocket.JSON.Send(client, msg) // send it to client
				check(err)
			}
		} else {
			err := websocket.JSON.Send(usernameConnection[msg.Reciptent], msg) // send it to individual reciptent
			check(err)
		}
	}
}

// Listens to the history channel and adds messages to the appropriate users history data
func addHistory() {

	for {
		// Listen on the channel and get new msgs
		channelMsg := <-historyChan

		// check if the msg needs timestamp and edit it
		msg := timeStamp(channelMsg)

		// If reciptent is forAll then send msg to all clients individual history
		if msg.Reciptent == "forAll" {
			for i, _ := range clientHistory {
				// Add message to each clients history
				if len(clientHistory[i]) < 1 { // Check if its the beggining of history and add message with no \n
					clientHistory[i] += msg.Message
				} else {
					clientHistory[i] += "\n" + msg.Message // add msg with new line
				}
			}
			// Add message to general history to be used by new client connecting who doesnt have individual history yet
			if len(mainHistory) < 1 { // check if its the beggining of history so dont add /n
				mainHistory += msg.Message
			} else {
				mainHistory += "\n" + msg.Message // add msg with new line
			}
		} else {
			// if message reciptent is an individual, send msg to the intended reciptent history only. Used for server msgs and errors
			if len(clientHistory[msg.Reciptent]) < 1 {
				clientHistory[msg.Reciptent] += msg.Message
			} else {
				clientHistory[msg.Reciptent] += "\n" + msg.Message
			}
		}
	}

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

// Reads welcome message from txt file to screen
func welcomeMsg() string {
	fileIO, _ := os.OpenFile("welcome.txt", os.O_RDWR, 0600) // open file
	defer fileIO.Close()
	rawBytes, _ := ioutil.ReadAll(fileIO) // read file
	return string(rawBytes) + "\n"
}

// Check if port is valid
func portChk(s string) bool {
	i, _ := strconv.Atoi(s)
	if i < 1024 || i > 65352 {
		log.Fatal("Your port number is not in the correct range of 1024-65352")
	}
	return true
}

// Sends user there chat history becuase after sending a message the client clears there screen
func sendHitstory (sockID *websocket.Conn, clientType string) {
	if clientType != "webapp" {
		msgHist := welcomeMsg() + "[ENTER YOUR NAME]: " + clientUserName[sockID] + "\n" + clientHistory[clientUserName[sockID]]
		sendHist := Message{Message: msgHist}
		time.Sleep(20 * time.Millisecond)			// DELETE/UDJUST if app screen refresh is not in sync
		err := websocket.JSON.Send(sockID, sendHist) // send history to client for when screen is refreshed
		check(err)
	}
}

// Check errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}