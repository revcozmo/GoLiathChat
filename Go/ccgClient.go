/************************

Go Command Chat
-Jeromy Johnson, Travis Lane
A command line chat system that 
will make it easy to set up a 
quick secure chat room for any 
number of people

************************/

package main

import (
	"./ccg"
	"./tboxgui"
	"container/list"
	"github.com/nsf/termbox-go"
	"time"
	"os"
)

type MessageObject struct {
	message   string //The Message
	sender    string //Who sent it
	timestamp int    //When it was sent
}

func main() {
	hostname := "jero.my:10234"

	if len(os.Args) > 1 {
		hostname = os.Args[1]
	}

	/* Initialize Connection */
	serv := ccg.NewHost()
	defer serv.Cleanup()
	err := serv.Connect(hostname)
	if err != nil {
		panic(err)
	}
	/* Initialize Termbox */
	termErr := termbox.Init()
	if termErr != nil {
		panic(termErr)
	}
	defer func() {
		termbox.Close()
	}()

	quit := false
	loggedin := false
	for !quit && !loggedin {
		quit, loggedin = displayLoginWindow(serv)
	}
	tboxgui.Clear()
	tboxgui.Flush()

	//loggedin,_ = serv.Login("Username", "password",0)
	if loggedin && !quit {
		// Start the server
		serv.Start()
		// Display the login window
		displayChatWindow(serv)
		quit = true
	}
}

// **********************************************
// **************** Login Window ****************
// ********************************************** 

func displayLoginWindow(serv *ccg.Host) (bool, bool) {
	quit := false
	login := false

	name := ""
	pass := ""
	login_err := ""
	// 0 Username 
	// 1 Password 
	// 2 Login
	// 3 Options
	// 4 Register
	box := 0

	const max_box = 4
	const min_box = 0

	//login_message := ""
	termboxEvent := make(chan termbox.Event)

	// Update the login window
	updateWindow := func() {
	
	
		tboxgui.Clear()
		sx, sy := termbox.Size()

		name_lines := tboxgui.GetLines(name, sx-2)
		pass_lines := tboxgui.GetLines(pass, sx-2)
		err_lines := tboxgui.GetLines(login_err, sx-2)

		tboxgui.WriteCenter((sy/2)-len(name_lines)-1, "Username:")
		tboxgui.WriteCenterWrap((sy/2)-len(name_lines), name_lines)
		tboxgui.WriteCenter((sy/2)+len(name_lines)+1, "Password:")
		tboxgui.WriteCenterWrap((sy/2)+len(name_lines)+2, pass_lines)

		by := (sy/2)+len(name_lines)+2+len(pass_lines)+1

		loginText := "Login"
		optionText := "Options"
		RegisterText := "Register"
		tboxgui.DrawButton(loginText, (box == 2), (sx/2)-(len(loginText) + 4) - ((len(optionText) +4)/2), by)
		tboxgui.DrawButton(optionText, (box == 3), (sx/2) - ((len(optionText) + 4) /2), by)
		tboxgui.DrawButton(RegisterText, (box == 4), 40, by)
		tboxgui.WriteCenterWrap(sy-len(err_lines), err_lines)
		tboxgui.Flush()
	}

	eventHandler := tboxgui.NewTermboxEventHandler()

	eventHandler.KeyEvents[termbox.KeyEnter] = func(_ termbox.Event) {
		if box == 0 {
			box = 1
		} else if box == 1 {
			if name == "" {
				login_err = "Username can not be blank."
			} else if pass == "" {
				login_err = "Password can not be blank."
			} else {
				login_err = "Logging in..."
				updateWindow()
				login, login_err = serv.Login(name, pass, 0)
			}
		} else if box == 2 {
			login_err = "Logging in..."
			updateWindow()
			login, login_err = serv.Login(name, pass, 0)
		} else if box == 3 {
			login_err = "Not implemented.."
			updateWindow()
		} else if box == 4 {
			displayRegisterWindow(serv, termboxEvent)
			updateWindow()
		}
	}
	eventHandler.KeyEvents[termbox.KeyCtrlC] = func(_ termbox.Event) {
		quit = true
		login = false
	}
	eventHandler.KeyEvents[termbox.KeyCtrlQ] = func(_ termbox.Event) {
		quit = true
		login = false
	}
	eventHandler.KeyEvents[termbox.KeyBackspace] = func(_ termbox.Event) {
		if box == 0 && len(name) > 0 { // Name
			name = name[0 : len(name)-1]
		} else if box == 1 && len(pass) > 0 { // Password
			pass = pass[0 : len(pass)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyBackspace2] = func(_ termbox.Event) {
		if box == 0 && len(name) > 0 { // Name
			name = name[0 : len(name)-1]
		} else if box == 1 && len(pass) > 0 { // Password
			pass = pass[0 : len(pass)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowUp] = func(_ termbox.Event) {
		if box > min_box {
			box -= 1
		} else {
			box = max_box
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowDown] = func(_ termbox.Event) {
		if box < max_box {
			box += 1
		} else {
			box = min_box
		}
	}
	eventHandler.KeyEvents[termbox.KeySpace] = func(_ termbox.Event) {
		if box == 0 && len(name) < 64 {
			name += " "
		} else if box == 1 && len(pass) < 64 {
			pass += " "
		}
	}
	eventHandler.KeyEvents[termbox.KeyTab] = func(_ termbox.Event) {
		if box == max_box {
			box = min_box
		} else {
			box += 1
		}
	}
	eventHandler.OnDefault = func(event termbox.Event) {
		if event.Ch != 0 {
			if box == 0 && len(name) < 64 {
				name += string(event.Ch)
			} else if box == 1 && len(pass) < 64 {
				pass += string(event.Ch)
			}
		}
	}

	updateWindow()

	// Start the goroutines
	go termboxEventPoller(termboxEvent)

	for !quit && !login {
		select {
		case event := <-termboxEvent:
			tboxgui.TermboxSwitch(event, eventHandler)
			updateWindow()
		}
	}

	updateWindow()

	return quit, login
}

// **********************************************
// ************** Register Window ***************
// ********************************************** 

func displayRegisterWindow(serv *ccg.Host, termboxEvent chan termbox.Event) {
	quit := false
	done := false
	username := ""
	password := ""
	passwordVerify := ""
	const max_box = 2
	const min_box = 0
	//  0 = Username
	//  1 = Password  
	//  2 = Password Verify
	box := min_box

	eventHandler := tboxgui.NewTermboxEventHandler()

	updateWindow := func() {
		tboxgui.Clear()
		sx, _ := termbox.Size()

		username_lines := tboxgui.GetLines(username, sx-2)
		password_lines := tboxgui.GetLines(password, sx-2)
		passwordVerify_lines := tboxgui.GetLines(passwordVerify, sx-2)

		tboxgui.WriteCenter(1, "Username")
		tboxgui.WriteCenterWrap(2, username_lines)
		tboxgui.WriteCenter(4, "Password")
		tboxgui.WriteCenterWrap(5, password_lines)
		tboxgui.WriteCenter(7, "Password Verify")
		tboxgui.WriteCenterWrap(8, passwordVerify_lines)

		tboxgui.Flush()
	}

	eventHandler.KeyEvents[termbox.KeyEnter] = func(_ termbox.Event) {
		if box == 0 {
			box = 1
		} else if box == 1 {
			box = 2
		} else if box == 2 {
			if password == passwordVerify {
				serv.Register(username, password)
			}
		}
	}
	eventHandler.KeyEvents[termbox.KeyCtrlC] = func(_ termbox.Event) {
		quit = true
		done = true
	}
	eventHandler.KeyEvents[termbox.KeyCtrlQ] = func(_ termbox.Event) {
		quit = true
		done = true
	}
	eventHandler.KeyEvents[termbox.KeyBackspace] = func(_ termbox.Event) {
		if box == 0 && len(username) > 0 { // Name
			username = username[0 : len(username)-1]
		} else if box == 1 && len(password) > 0 { // Password
			password = password[0 : len(password)-1]
		} else if box == 1 && len(passwordVerify) > 0 { // Password
			passwordVerify = passwordVerify[0 : len(passwordVerify)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyBackspace2] = func(_ termbox.Event) {
		if box == 0 && len(username) > 0 { // Name
			username = username[0 : len(username)-1]
		} else if box == 1 && len(password) > 0 { // Password
			password = password[0 : len(password)-1]
		} else if box == 1 && len(passwordVerify) > 0 { // Password
			passwordVerify = passwordVerify[0 : len(passwordVerify)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowUp] = func(_ termbox.Event) {
		if box > min_box {
			box -= 1
		} else {
			box = max_box
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowDown] = func(_ termbox.Event) {
		if box < max_box {
			box += 1
		} else {
			box = min_box
		}
	}
	eventHandler.KeyEvents[termbox.KeySpace] = func(_ termbox.Event) {
		if box == 1 && len(password) < 64 {
			password += " "
		} else if box == 2 && len(passwordVerify) < 64 {
			passwordVerify += " "
		}
	}
	eventHandler.OnDefault = func(event termbox.Event) {
		if event.Ch != 0 {
			if box == 0 && len(username) < 64 {
				username += string(event.Ch)
			} else if box == 1 && len(password) < 64 {
				password += string(event.Ch)
			} else if box == 2 && len(passwordVerify) < 64 {
				passwordVerify += string(event.Ch)
			}
		}
	}

	updateWindow()
	for !quit && !done {
		select {
		case event := <-termboxEvent:
			tboxgui.TermboxSwitch(event, eventHandler)
			updateWindow()
		}
	}
}

// *********************************************
// **************** Chat Window ****************
// ********************************************* 

func displayChatWindow(serv *ccg.Host) {

	// Setup the variables
	input := ""
	running := true
	start_message := 0
	messages := list.New()
	termboxEvent := make(chan termbox.Event)
	eventHandler := tboxgui.NewTermboxEventHandler()

	eventHandler.KeyEvents[termbox.KeyCtrlQ] = func(_ termbox.Event) {
		tboxgui.Clear()
		tboxgui.MessageUs("Exiting...")
		tboxgui.Flush()
		time.Sleep(time.Second * 2)
		running = false
	}
	eventHandler.KeyEvents[termbox.KeyCtrlC] = func(_ termbox.Event) {
		tboxgui.Clear()
		tboxgui.Flush()
		running = false

	}
	eventHandler.KeyEvents[termbox.KeyEnter] = func(_ termbox.Event) {
		if input != "" {
			serv.Send(input)
			input = ""
		}
	}
	eventHandler.KeyEvents[termbox.KeyBackspace] = func(_ termbox.Event) {
		if len(input) > 0 {
			input = input[0 : len(input)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyBackspace2] = func(_ termbox.Event) {
		if len(input) > 0 {
			input = input[0 : len(input)-1]
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowUp] = func(_ termbox.Event) {
		if start_message < messages.Len() {
			start_message += 1
		}
	}
	eventHandler.KeyEvents[termbox.KeyArrowDown] = func(_ termbox.Event) {
		if start_message > 0 {
			start_message -= 1
		}
	}
	eventHandler.KeyEvents[termbox.KeySpace] = func(_ termbox.Event) {
		if len(input) <= 160 {
			input += " "
		}
	}
	eventHandler.OnDefault = func(keyEvent termbox.Event) {
		if keyEvent.Ch != 0 {
			if len(input) <= 160 {
				input += string(keyEvent.Ch)
			}
		}
	}
	eventHandler.OnResize = func(_ termbox.Event) {

	}

	// Display the window
	tboxgui.Clear()
	updateChatWindow(input, messages, start_message)
	tboxgui.Flush()
	// Start the goroutines
	go termboxEventPoller(termboxEvent)
	// Run the main loop
	for running {
		select {
		case event := <-termboxEvent:
			tboxgui.TermboxSwitch(event, eventHandler)
			tboxgui.Clear()
			updateChatWindow(input, messages, start_message)
			tboxgui.Flush()
		case serverEvent := <-serv.Reader:
			message := MessageObject{string(serverEvent.Payload), serverEvent.Username, time.Now().Second()}
			messages.PushFront(message)
			tboxgui.Clear()
			updateChatWindow(input, messages, start_message)
			tboxgui.Flush()
		}
	}
}

//Updates the chat
func updateChatWindow(input string, messages *list.List, start_message int) {

	x, y := termbox.Size()
	if x != 0 && y != 0 {
		input_top := displayInput(input)
		displayMessages(messages, start_message, input_top)
	}
}

// Displays the chat messages
func displayMessages(messages *list.List, offset int, input_top int) {
	line_cursor := input_top
	sx, sy := termbox.Size()
	// Iterate to the current message
	p := messages.Front()
	for i := 0; i < offset; i++ {
		p = p.Next()
	}
	// Iterate over the messages
	for ; p != nil; p = p.Next() {

		cur := p.Value.(MessageObject)
		lines := tboxgui.GetLines(cur.message, sx)
		//tboxgui.FillH("-", 0, sy-line_cursor, sx)

		line_cursor += 1
		for i := len(lines) - 1; i >= 0; i-- {
			tboxgui.Write(0, sy-line_cursor, lines[i])
			line_cursor += 1
		}
		if p.Next() != nil {
			if p.Next().Value.(MessageObject).sender == cur.sender {
				line_cursor -= 1
			} else {
				tboxgui.Write(0, sy-line_cursor, cur.sender)
				tboxgui.FillH("-", len(cur.sender), sy-line_cursor, sx)
			}
		} else {
			tboxgui.Write(0, sy-line_cursor, cur.sender)
			tboxgui.FillH("-", len(cur.sender), sy-line_cursor, sx)
		}
	}
}

// Displays the chat input
func displayInput(input string) int {
	sx, sy := termbox.Size()
	line_cursor := 1
	if input == "" {
		termbox.SetCursor(0, sy-line_cursor)
		line_cursor += 1
		tboxgui.FillH("-", 0, sy-line_cursor, sx)
		return line_cursor
	} else {
		lines := tboxgui.GetLines(input, sx)
		for i := len(lines) - 1; i >= 0; i-- {
			tboxgui.Write(0, sy-line_cursor, lines[i])
			line_cursor += 1
		}
		termbox.SetCursor(len(lines[len(lines)-1]), sy-1)
		tboxgui.FillH("-", 0, sy-line_cursor, sx)
		return line_cursor
	}
	return 1
}

// *********************************************
// ************** Other Functions **************
// ********************************************* 

func termboxEventPoller(event chan<- termbox.Event) {
	for {
		event <- termbox.PollEvent()
	}
}
