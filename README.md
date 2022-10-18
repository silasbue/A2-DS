# A2-DS

Hand-in 2 for Distributed Systems course, ITU SWU 3.sem

[Link to Assignment](https://learnit.itu.dk/mod/assign/view.php?id=165135)

## a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?

We have made a implementation of a package, which is represented by a struct called '_package'. This consist of meta-data (syncronization and acknowledgement keys) and data (the data-field, which consist of an int).

## b) Does your implementation use threads or processes? Why is it not realistic to use threads?

Our implementation uses go routines (threads) to simulate the process. This is not realistic in a real implementation since the process would happen between different systems across a network. Using goroutines is not a realistic simulation since they are much more stable, than client and servers using the network. Fx, in the real world, a server or client can fail (turn off, crash, etc), where as our threads won't.

## c) How do you handle message re-ordering?

Re-ordering is when the messages don't arrive in the same order, which they were sent. This can cause problems in multiple cases. One example is if a web server sends an html file to a client. If this message is split up into fragments, reordered and then assembled in an incorrect order, this would lead to big problems.
When a package is sent in multiple fragments, we sort the fragments by their sequence number, thereby unsuring that the data is interpreted in the correct order.

## d) How do you handle message loss?

Message loss is when a message is sent, but never gets received. The problems with this are obvious. This can be caused by multiple things. Device failures, connection loss (physical and wireless), etc. Here it is important that a system won't wait forever until a message arrives, since if the message is lost, the system will halt forever.
We handle message loss by keep sending a package until we get a response. This ensures that the other part has "heard" the message.
Though there is one case, where a message loss will break the program. This is in the case that the ACK message is lost, when it should be a response to a SYN+ACK message.

## e) Why is the 3-way handshake important?

To make sure that the connection is established and to assure that the receiver knows which packages to receive, and also knows when a package has been lost.
