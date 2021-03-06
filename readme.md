#GoLiath Chat

	So small that it's a big deal!
	
GoLiath chat is a simple secure server-client chat system with 
login and verification implemented in go.

Currently we are using a custom webkit browser and html/js/css for our UI. We chose this method as it is very cross platform (everyone can open a webpage, right?), and allows for lots of customization. Also, it is very lightweight compared to other UI Kit frameworks (I'm looking at you .NET). 

We use port 10234 to communicate, be sure to have it open if you are running a server. 

##Goals
Our main goal for this project is to have a quick and secure chatroom that will be as simple as possible to set up and connect to. 

##Server
Features:
- Retrievable chat history
- Secure communication
- File transfers
- Small footprint

##Client
**Login**

The first screen that comes up will be the login screen, enter your username and password to authenticate with the server.
If you do not have a username and password on the server, click register to bring up the registration screen.

**Registration**

From the registration screen, enter your desired username, password and server, then submit it. Your request gets sent to the server where it is reviewed by moderators and accepted or denied. If it is accepted you may then log into the server with your username and password.
Note: Registrations are per-server and your account details on one server will not get you into another server.

**History**

Upon login, the server will send a large chunk of chat history to you.

**File Transfers**

To upload a file to the server for sharing, use:

	/upload path/to/filename

Once a file is on the server, other users may download it with:

	/dl filename

##Installation (Building from source)

Installation on Linux is very simple, ensure you have Go installed and do the following:

	git clone git://github.com/whyrusleeping/GoLiathChat.git
	cd GoLiathChat
	./configure.sh

And depending on which version you want (Goliath vs GoliathLite)

	#For Goliath (requires gtk webkit libraries)
	make GoliathWK

	#For GoliathLite (statically builds with no need for external libraries)
	make GoliathLite
