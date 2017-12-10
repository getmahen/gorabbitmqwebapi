# Welcome to Buffalo!

Thank you for choosing Buffalo for your web development needs.


## Database Setup

It looks like you chose to set up your application using a postgres database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start postgres for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

	$ buffalo db create -a


## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

## What Next?

We recommend you heading over to [http://gobuffalo.io](http://gobuffalo.io) and reviewing all of the great documentation there.

Good luck!

[Powered by Buffalo](http://gobuffalo.io)

## WHAT IS THIS APP DO?
On Start Up:
This Web API 
- Initializes the routes
- Kick starts RabbitMq with a queue and with default settings
- Fires up a GO ROUTINE to start listening on the Queue messages. This runs in the background and starts consuming messages from the queue

This is a Web API that exposes 2 Endpoints 
 1. /PortRequest - This endpoint sends an outgoing request to validate given phone number's Port eligibility (Note: the out going call to an external API is not implemented). It also creates a unique RequestID for every request and initializes a GO Channel and starts waiting for a message on the Channel (this causes a block in the request execution process). Ofcourse, there is a Timeout as well on the channel to create a short circuit.

 2. /PortCallback - This is a callback URL that accepts response from the third part API (PortEligibility Asynchronous response) for a given RequestId. It then puts the message on the RabbitMq queue. One the message for a RequestId is received on the Queue (broker.beginConsumption() reaps the message from the RabbitMq's Queue and calls the listener), it sends the same message to Channel. This unblocks the WAIT that is initialized by the  /PortRequest endpoint