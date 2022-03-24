# Server Setup Instructions:

In order to start the chat server type: **go run .** <br/>
If no arguments are passed the default port is 8989. <br/>
Optionally you can set your own port: **go run . 4444** will start a server on port 4444 <br/>


# Chat Client Setup Instructions:

In order to start the client app. <br/>
Open a terminal session in the chat-client directory. <br/>
Install it using the command: ** ** <br/>
You can run the chat client with default settings by typing **"nc"** <br/>
This would use server address: localhost and port: 8989. <br/>
In order to set custom address and port, use syntax **nc "server-address" "port"** <br/>
Example: **"nc localhost 8989"** would start the chat client connecting to chat server on localhost at port 8989.

# Web Front End
The application also has a simple web front end which allows users to interact wether in terminal or on the website. <br/>
To access simply go to http://localhost:8989/<br/>
If you would like to allow access to remote users then go to /public/app.js and replace all mentions of localhost with your public ip address. <br/>
Users can then access the web front end at (your local ip):8989/

# Change Username:

In order to change your chat username. Simply send a message saying: <br/>
**"hey server change my name to"** and add your new name. <br/>
The server will respond and tell you it has changed your name.

# Time.Sleep

The app is designed so that the client clears the screen immediately after sending a message. <br/>
The server knows this and immediately sends a copy of that clients chat history to there screen. <br/>
However, depending on how fast your computer is to clear the screen. The timing of the message sent by the server can be off. <br/>
To circumvent this there is a time delay in the servers response. If your app is taking to long to get chat history and you have a blank screen after sending a message. You can change the delay in the code, or remove it completely. <br/>
Search in netCat.go for "**time.Sleep(20 * time.Millisecond)**" and adjust accordingly.