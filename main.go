package main

import (
	"fmt"
	"os"
	"time"
)

var requestDone bool = false

func main() {
	clientAmount := 5

	// Initialize channels
	forwarderChann := make(chan _package)
	var clientChannels []chan _package
	for i := 0; i < clientAmount; i++ {
		clientChannels = append(clientChannels, make(chan _package))
	}

	var userChannels []chan userRequest
	for i := 0; i < clientAmount; i++ {
		userChannels = append(userChannels, make(chan userRequest))
	}

	// Run goroutines
	go postNord(forwarderChann, clientChannels)

	for i, chann := range clientChannels {
		go client(i, chann, forwarderChann, userChannels[i])
	}

	// Run user interface
	userInterface(userChannels)

	// What to simulate (perfect, message loss, re-ordering, etc)

	// Which client to communicate between

}

func client(ip int, inbox chan _package, outgoing chan _package, userInbox chan userRequest) {
	sendPackage := func(p _package) {
		outgoing <- p
	}
	receivePackage := func(p _package) {
		p = <-inbox
	}

	exchangePackages := func(p _package) _package {
		recieved := false

		go func() {
			for !recieved {
				rp := <-inbox
				if p.seq != 0 && p.ack == 0 && p.fin == 0 && p.data == "" { //
					recieved = rp.seq != 0 && rp.ack != 0 && p.fin == 0 && p.data == ""
				} else if p.seq != 0 && p.ack != 0 && p.fin == 0 && p.data == "" {
					recieved = rp.seq == 0 && rp.ack != 0 && p.fin == 0 && p.data == ""
				} else if p.seq == 0 && p.ack == 0 && p.fin == 0 && p.data != "" {
					recieved = rp.seq == 0 && rp.ack != 0 && rp.fin == 0 && rp.data != ""
				} else if p.seq == 0 && p.ack != 0 && rp.fin == 0 && p.data != "" {
					recieved = rp.seq == 0 && rp.ack != 0 && rp.fin == 0 && rp.data != ""
				}
			}
		}()

		for !recieved {
			outgoing <- p
			time.Sleep(10 * time.Millisecond)
		}

		return p
	}

	verifyAck := func(SYN, ACK int) bool {
		return SYN+1 == ACK
	}

	for {
		select {
		case request := <-userInbox: // Sender

			// Three-way handshake
			p := exchangePackages(_package{ip, request.desIP, 1, 0, 0, ""})

			outgoing <- _package{ip, request.desIP, 0, p.seq + 1, 0, ""}

			message := request.data

			for i := 0; i < len(message); i++ {
				seq++
				outgoing <- _package{ip, request.desIP, 0, seq, 0, string(message[i])}
			}

		case p := <-inbox: // Receiver

			if p.seq != 0 && p.ack == 0 {
				p = exchangePackages(_package{ip, p.srcIP, 1, p.seq + 1, 0, ""})
			}
			if p.seq == 0 && p.ack != 0 {

			}

		}
	}

}

func postNord(chann chan _package, clientChannels []chan _package) {
	for {
		// Listen for packages
		p := <-chann

		// Send package to client
		// Here, message loss, re-ordering, delay, etc will be simulated
		clientChannels[p.desIP] <- p
	}
}

type _package struct {
	srcIP int
	desIP int
	seq   int
	ack   int
	fin   int
	data  string
}

type userRequest struct {
	desIP    int
	data     string
	simError int
}

// userInterface for managing the program
func userInterface(chann []chan userRequest) {
	fmt.Println("Welcome to a TCP/IP simulation")
	fmt.Println("------------------------------")

	for {
		requestDone = false
		fmt.Println("Please enter a value to start")
		fmt.Println("1) Start a new simulation")
		fmt.Println("2) Exit")

		var input int
		fmt.Scanln(&input)

		if input == 2 {
			os.Exit(3)
		}

		fmt.Println("Please choose a client to communicate from")
		fmt.Println("0) Thore")
		fmt.Println("1) Riko")
		fmt.Println("2) Troels")
		fmt.Println("3) Rasmus")

		var from int
		fmt.Scanln(&from)

		fmt.Println("Please choose a client to communicate to")
		if from == 0 {
			fmt.Println("1) Riko")
			fmt.Println("2) Troels")
			fmt.Println("3) Rasmus")
		} else if from == 1 {
			fmt.Println("0) Thore")
			fmt.Println("2) Troels")
			fmt.Println("3) Rasmus")
		} else if from == 2 {
			fmt.Println("0) Thore")
			fmt.Println("1) Riko")
			fmt.Println("3) Rasmus")
		} else if from == 3 {
			fmt.Println("0) Thore")
			fmt.Println("1) Riko")
			fmt.Println("2) Troels")
		}

		var to int
		fmt.Scanln(&to)

		fmt.Println("Please write a message to send:")
		var message string
		fmt.Scanln(&message)

		fmt.Println("Choose simulation cases:")
		fmt.Println("0) Perfect")
		fmt.Println("1) Message loss")
		fmt.Println("2) Re-ordering")
		fmt.Println("3) Delay")
		fmt.Println("4) All")

		var simulationCase int
		fmt.Scanln(&simulationCase)

		fmt.Println("Starting simulation...")
		fmt.Println("------------------------------")

		newRequest := userRequest{to, message, simulationCase}
		chann[from] <- newRequest

		for !requestDone {
			time.Sleep(1 * time.Second)
		}
	}

}
