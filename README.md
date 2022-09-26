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

We dont handle message re-ordering, since we only send one message at a time.

## d) How do you handle message loss?

We dont handle message loss, since we have simulated this part of the process.

## e) Why is the 3-way handshake important?

To make sure that the connection is established and that to assure that the receiver knows which packages to receive, and also knows when a package has been lost.
