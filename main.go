package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	SYNBIT   = 1 << 0
	ACKBIT   = 1 << 1
	FINBIT   = 1 << 2
	DATABIT  = 1 << 3
	STARTBIT = 1 << 4
)

func isAppropriateResponse(sent _package, response _package) bool {
	if sent.dataOffset == 1 { //			Case	(SYN)
		return response.dataOffset == 3 // 	Response	(SYN+ACK)
	} else if sent.dataOffset == 3 { //		Case	(SYN+ACK)
		return response.dataOffset == 2 //	Response	(ACK)
	} else if sent.dataOffset == 8 { //		Case	(DATA)
		return response.dataOffset == 10 //	Response	(DATA+ACK)
	} else if sent.dataOffset == 16 { //	Case	(START)
		return response.dataOffset == 1 //	Response	(SYN)
	} else if sent.dataOffset == 4 { //		Case	(FIN)
		return response.dataOffset == 6 //	Response	(FIN+ACK)
	} else if sent.dataOffset == 6 { //		Case	(FIN+ACK)
		return response.dataOffset == 2 //	Response	(ACK)
	}
	return false
}

var requestDone bool = false

func main() {
	rand.Seed(1)
	clientAmount := 4

	// Initialize channels
	forwarderChann := make(chan _package)
	var clientChannels []chan _package
	simulationCaseChann := make(chan int)
	for i := 0; i < clientAmount; i++ {
		clientChannels = append(clientChannels, make(chan _package, 8))
	}

	// Run goroutines for clients
	for i, chann := range clientChannels {
		go client(i, chann, forwarderChann)
	}

	go forwarder(forwarderChann, clientChannels, simulationCaseChann)

	// Run user interface
	userInterface(forwarderChann, simulationCaseChann)

	// What to simulate (perfect, message loss, re-ordering, etc)

	// Which client to communicate between
}

func client(ip int, inbox chan _package, outgoing chan _package) {
	SEQ := 0
	ACK := 0
	FIN := 0
	var dataPackagesReceived []_package
	var dataToSend []string
	dataToSendIdx := 0
	sendPackage := func(p _package) {
		outgoing <- p
		fmt.Println("IP" + fmt.Sprint(ip) + ": Sent package: " + fmt.Sprint(p))
	}
	receivePackage := func() _package {
		p := <-inbox
		fmt.Println("IP" + fmt.Sprint(ip) + ": Received package: " + fmt.Sprint(p))
		return p
	}

	exchangePackages := func(p _package) _package {
		var receivedPackage _package

		go func() {
			for (receivedPackage == _package{}) || !isAppropriateResponse(p, receivedPackage) {
				receivedPackage = receivePackage()
			}
		}()
		for (receivedPackage == _package{}) || !isAppropriateResponse(p, receivedPackage) {
			sendPackage(p)
			time.Sleep(10 * time.Millisecond)
		}

		return receivedPackage
	}

	verifyAck := func(SYN, ACK int) bool {
		return SYN+1 == ACK
	}

	for {
		rp := receivePackage()
		if rp.dataOffset == 16 { // Case 0 (Start)

			dataToSend = strings.Split(rp.data, " ")

			SEQ = rand.Int()
			var p _package
			for (p == _package{}) || !verifyAck(SEQ, p.ack) {
				p = exchangePackages(_package{ip, rp.srcIP, SYNBIT, SEQ, 0, 0, ""})
			}
			inbox <- p
			//fmt.Println("\tCase 0 completed")
		} else if rp.dataOffset == 1 { // Case 1 (SYN)
			SEQ = rand.Int()
			ACK = rp.seq + 1
			var p _package
			for (p == _package{}) || !verifyAck(SEQ, p.ack) {
				p = exchangePackages(_package{ip, rp.srcIP, SYNBIT | ACKBIT, SEQ, ACK, 0, ""})
			}
		} else if rp.dataOffset == 3 { // Case 2 (SYN+ACK)
			ACK = rp.seq + 1
			sendPackage(_package{ip, rp.srcIP, ACKBIT, 0, ACK, 0, ""})
			SEQ++
			inbox <- exchangePackages(_package{ip, rp.srcIP, DATABIT, SEQ, 0, 0, dataToSend[dataToSendIdx]})
			dataToSendIdx++
		} else if rp.dataOffset == 8 { // Case 3 (DATA)
			if rp.seq == ACK {
				dataPackagesReceived = append(dataPackagesReceived, rp)
				ACK++
			}
			sendPackage(_package{ip, rp.srcIP, DATABIT | ACKBIT, 0, rp.seq + 1, 0, ""})
		} else if rp.dataOffset == 10 { // Case 4 + 5 (DATA+ACK)
			if dataToSendIdx < len(dataToSend) { // Case 4 (DATA+ACK) - send more data
				SEQ++
				inbox <- exchangePackages(_package{ip, rp.srcIP, DATABIT, SEQ, rp.seq + 1, 0, dataToSend[dataToSendIdx]})
				dataToSendIdx++
				ACK++
			} else { // Case 5 (DATA+ACK) - send FIN
				FIN = rand.Int()
				p := exchangePackages(_package{ip, rp.srcIP, FINBIT, 0, 0, FIN, ""})
				ACK = p.fin + 1
				inbox <- p
			}
		} else if rp.dataOffset == 4 { // Case 6 (FIN)
			// Sort packages (message re-ordering)
			sort.Slice(dataPackagesReceived, func(i, j int) bool { return dataPackagesReceived[i].seq < dataPackagesReceived[j].seq })

			FIN = rand.Int()
			inbox <- exchangePackages(_package{ip, rp.srcIP, FINBIT | ACKBIT, 0, ACK, FIN, ""})
			// Print data from packages
			fmt.Print("IP" + fmt.Sprint(ip) + " has received all data packages: ")
			for _, p := range dataPackagesReceived {
				fmt.Print(p.data)
				fmt.Print(" ")
			}
			fmt.Println()
		} else if rp.dataOffset == 6 { // Case 7 (FIN+ACK)
			sendPackage(_package{ip, rp.srcIP, ACKBIT, 0, rp.fin + 1, FIN, ""})
			FIN++
			requestDone = true
		}
	}
}

func forwarder(chann chan _package, clientChannels []chan _package, simulationCaseChann chan int) {
	simulationCase := <-simulationCaseChann
	fmt.Println("Chosen simulation case: " + fmt.Sprint(simulationCase))
	for {
		// Listen for packages
		p := <-chann

		// Send package to client
		// Here, message loss, re-ordering, delay, etc will be simulated
		if simulationCase == 0 || p.dataOffset == 16 {
			clientChannels[p.desIP] <- p
		} else if simulationCase == 1 {
			action := rand.Float32()
			if action < 0.2 {
				fmt.Println("Ups, message lost in transport")
			} else {
				clientChannels[p.desIP] <- p
			}
		} else {
			panic("Unknown simulation case!")
		}
	}
}

type _package struct {
	srcIP      int
	desIP      int
	dataOffset uint8
	// 0000 0000
	//    ^ ^^^^
	//Start ||||
	//   Data|||
	//     FIN||
	//      ACK|
	//       SYN

	seq  int
	ack  int
	fin  int
	data string
}

type userRequest struct {
	desIP    int
	data     string
	simError int
}

// userInterface for managing the program
func userInterface(forwarder chan _package, simulationChann chan int) {
	fmt.Println("Welcome to a TCP/IP simulation")
	fmt.Println("------------------------------")
	clientNames := []string{"Thore", "Riko", "Troels", "Rasmus"}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter a value to start")
	fmt.Println("1) Start a new simulation")
	fmt.Println("2) Exit")

	input, _ := reader.ReadString('\n')
	input = strings.Trim(input, "\n")
	if input == "2" {
		os.Exit(3)
	}

	fmt.Println("Please choose a client to communicate from")
	for i, name := range clientNames {
		fmt.Println(fmt.Sprint(i) + ") " + name)
	}

	from, _ := reader.ReadString('\n')
	from = strings.Trim(from, "\n")
	fromIdx, _ := strconv.Atoi(from)
	fmt.Println("Sending from: " + fmt.Sprint(fromIdx))

	fmt.Println("Please choose a client to communicate to")

	for i, name := range clientNames {
		if i == fromIdx {
			continue
		}
		fmt.Println(fmt.Sprint(i) + ") " + name)
	}

	to, _ := reader.ReadString('\n')
	to = strings.Trim(to, "\n")
	toIdx, _ := strconv.Atoi(to)

	fmt.Println("Please write a message to send:")
	message, _ := reader.ReadString('\n')
	message = strings.Trim(message, "\n")

	fmt.Println("Choose simulation cases: (NOT IMPLEMENTED!)")
	fmt.Println("0) Perfect")
	fmt.Println("1) Message loss")

	simulationCase, _ := reader.ReadString('\n')
	simulationCase = strings.Trim(simulationCase, "\n")
	simulationCaseIdx, _ := strconv.Atoi(simulationCase)

	simulationChann <- simulationCaseIdx

	fmt.Println("Starting simulation...")
	fmt.Println("------------------------------")

	newRequest := _package{toIdx, fromIdx, STARTBIT, 0, 0, 0, message}
	forwarder <- newRequest

	for !requestDone {
		time.Sleep(100 * time.Millisecond)
	}
}
