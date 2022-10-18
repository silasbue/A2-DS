package main

import (
	"fmt"
	"sort"
	"time"
)

func main() {
	inClient1 := make(chan _package)
	outClient1 := make(chan _package)

	go client(1, inClient1, outClient1)

	// establishing connection using three way handshake
	inClient1 <- _package{1, 0, "", "syn"}
	<-outClient1
	inClient1 <- _package{2, 2, "", "ack"}

	// sending data
	inClient1 <- _package{3, 0, "test", "data"}
	<-outClient1

	// sending packages in the wrong order to test the reordering
	inClient1 <- _package{5, 0, "test3", "data"}
	<-outClient1
	inClient1 <- _package{4, 0, "test2", "data"}
	<-outClient1

	// closing connection
	inClient1 <- _package{5, 0, "", "fin"}
	<-outClient1
	inClient1 <- _package{4, 7, "", "ack"}

	time.Sleep(2 * time.Second)
}

func client(id int, inbox chan _package, outgoing chan _package) {

	inboxHandler := func() {
		thisSeq := 1
		var recievedData []_package
		// ingoingConnectionEstablished := false
		for {
			recievedMsg := <-inbox
			switch recievedMsg.flag {
			case "syn":
				fmt.Println("Client", id, "syn recieved")
				outgoing <- _package{thisSeq, recievedMsg.seq + 1, "", "syn+ack"}
				fmt.Println("Client", id, "syn+ack sent")
				break
			case "syn+ack":
				if recievedMsg.ack != thisSeq {
					fmt.Println("Client", id, "recieved WRONG syn+ack")
				} else {
					fmt.Println("Client", id, "recieved CORRECT syn+ack")
				}
				outgoing <- _package{thisSeq, recievedMsg.seq + 1, "", "ack"}
				break
			case "ack":
				if recievedMsg.ack != thisSeq {
					fmt.Println("Client", id, "recieved WRONG ack")
					// ingoingConnectionEstablished = true

				} else {
					fmt.Println("Client", id, "recieved CORRECT ack")
					fmt.Println("Client", id, "ingoing connection established")
				}
				break
			case "data":
				outgoing <- _package{thisSeq, recievedMsg.seq + 1, "", "ack"}
				recievedData = append(recievedData, recievedMsg)
				fmt.Println("Client", id, "recieved the following data: ", recievedMsg.data)
				break
			case "fin":
				fmt.Println("Client", id, "fin recieved")
				outgoing <- _package{thisSeq, recievedMsg.seq + 1, "", "fin+ack"}
				fmt.Println("Client", id, "fin+ack sent")
				// ingoingConnectionEstablished = false
				fmt.Println("Client", id, "has closed an ingoing connection")
				//reorder data
				sort.Slice(recievedData, func(i, j int) bool {
					return recievedData[i].seq < recievedData[j].seq
				})
				fmt.Println(recievedData)
				break
			case "fin+ack":
				if recievedMsg.ack != thisSeq {
					fmt.Println("Client", id, "recieved WRONG fin+ack")
				} else {
					fmt.Println("Client", id, "recieved CORRECT fin+ack")
				}
				outgoing <- _package{thisSeq, recievedMsg.seq + 1, "", "ack"}
				fmt.Println("Client", id, "has closed an outgoing connection")
				break
			}
			// increment the sequence number for the next message
			thisSeq++
		}
	}
	go inboxHandler()

	time.Sleep(2 * time.Second)
}

type _package struct {
	seq  int
	ack  int
	data string
	flag string
}
