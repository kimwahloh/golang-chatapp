# GoLang Chat Application

This repository contains a GoLang-based chat application that allows users to communicate via chat rooms and direct messages. Users can run both the server and client components to connect to the chat system.

## Usage Guidelines

### (1) Running the Application

#### Server:

Navigate to the `server` folder.

Run the server using:
go run .


#### Client:

Navigate to the `client` folder.

Run the client using:
go run .


### (2) Connecting to Public Chat Room

Upon logging in, users are automatically connected to the public chat room.

Type `/` for available commands.

### (3) Available Commands

- `/list`: Displays a list of available chat rooms.
- `/join <room_name>`: Joins the specified chat room.
- `/leave <room_name>`: Leaves the specified chat room.
- `/create <room_name>`: Creates a new chat room with the specified name.
- `/dm <recipient> <message>`: Sends a direct message to the specified recipient.
- `/save`: Save message logs of the user.
