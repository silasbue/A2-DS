# A2-DS

Hand-in 2 for Distributed Systems course, ITU SWU 3.sem

[Link to Assignment](https://learnit.itu.dk/mod/assign/view.php?id=165135)

## a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?

We have made a implementation of a package, which is represented by a struct called 'present'. This consist of meta-data (syncronization and acknowledgement keys) and data (the data-field, which consist of an int).

## b) Does your implementation use threads or processes? Why is it not realistic to use threads?

yes, our implementation uses go routines to simulate the process.
This is not realistic in a real implementation since the process would happen between different systems across a network.
Threads does not make sense in

## c) How do you handle message re-ordering?

Re-ordering is when the messages don't arrive in the same order, which they were sent. This can cause problems in multiple cases. One example is if a web server sends an html file to a client. If this message is split up into fragments, reordered and then assembled in an incorrect order, this would lead to big problems.
We dont handle message re-ordering, since we only send one message at a time. This message also won't be split into fragments, which therefore can't be re-ordered.

## d) How do you handle message loss?

Message loss is when a message is sent, but never gets received. The problems with this are obvious. This can be caused by multiple things. Device failures, connection loss (physical and wireless), etc. Here it is important that a system won't wait forever until a message arrives, since if the message is lost, the system will halt forever.
We dont handle message loss, since we have simulated this part of the process.

## e) Why is the 3-way handshake important?

To make sure that the connection is established and to assure that the receiver knows which packages to receive, and also knows when a package has been lost.
