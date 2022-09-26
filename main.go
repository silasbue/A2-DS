package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan present)

	go client(ch)
	go server(ch)

	time.Sleep(10 * time.Second)
}

func client(chann chan present) {
	//Inital sending containing sync-key
	clientsync := 1
	sending := present{clientsync, 0, 0}

	chann <- sending
	fmt.Println("Client: syn sent")

	//Recieves ack-key and verifies it
	received := <-chann
	ack := received.ack
	time.Sleep(10 * time.Millisecond)

	if ack == clientsync+1 {
		fmt.Println("Client: correct syn+ack recieved")
	}

	serverSync := received.sync

	//Sends ack-key to server
	sending = present{0, serverSync + 1, 0}
	chann <- sending
	fmt.Println("Client: ack sent")

	//Sends data to server
	sending = present{0, 0, 9203947783953}
	chann <- sending
	fmt.Println("Client: data package sent")
}

func server(chann chan present) {
	//Recieves syn-key
	received := <-chann
	time.Sleep(10 * time.Millisecond)
	clientSync := received.sync
	serverSync := 4
	fmt.Println("Server: syn recieved")

	//Sends syn+ack-key to client
	sending := present{serverSync, clientSync + 1, 0}
	chann <- sending
	fmt.Println("Server: syn+ack sent")

	//Recieves ack-key and verifies it
	received = <-chann
	time.Sleep(10 * time.Millisecond)

	if received.ack == serverSync+1 {
		fmt.Println("Server: correct ack recieved")
		//Recieves data from client
		received = <-chann
		time.Sleep(10 * time.Millisecond)
		fmt.Println("Server: data recieved, is: ", received.data)
	} else {
		fmt.Println("Server: Recieved wrong ack")
	}
}

type present struct {
	sync int
	ack  int
	data int
}
