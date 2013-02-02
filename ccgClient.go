package ccg

import (
	"github.com/nsf/termbox-go"
	"time"
	"net"
	"bytes"
	"encoding/binary"
)

func main() {
	defer cleanup()
	hostname := "127.0.0.1"
	messages := make(chan Packet)
	writer, err := makeConnection(hostname,messages)
	if err != nil {
		panic(err)
	}
	//send a test packet, so tired right now. Its probably all broken
	p := Packet{}
	p.typ = 1
	p.payload = "hello world"
	writer <- p
	ui()
}

// Network
func makeConnection(hostname string, mesChan chan<- Packet) (chan<- Packet, error) {
	addr, err := net.ResolveTCPAddr("tcp",hostname)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp",nil,addr)
	if err != nil {
		return nil, err
	}
	writer := make(chan Packet)
	go writeMessages(conn, writer)
	go readMessages(conn, mesChan)
	return writer, nil
}

func writeMessages(conn *net.TCPConn, writeChan <-chan Packet) {

}

func readMessages(conn *net.TCPConn, mesChan chan<- Packet) {
	flagBuf := make([]byte, 1)
	lenBuf  := make([]byte, 2)
	timeBuf := make([]byte, 4)
	for {
		//Need to check connectivity to see if a disconnect has happened
		p := Packet{}
		conn.Read(flagBuf)
		p.typ = flagBuf[0] //Packet is just one byte
		conn.Read(timeBuf)
		buf := bytes.NewBuffer(timeBuf)
		binary.Read(buf, binary.LittleEndian, &p.timestamp)
		conn.Read(lenBuf)
		buf = bytes.NewBuffer(lenBuf)
		binary.Read(buf, binary.LittleEndian, &p.mesLen)
		strBuf := make([]byte, p.mesLen)
		conn.Read(strBuf)
		p.payload = string(strBuf)
		mesChan <- p
	}
}



// Handles login functions, returns true (successful) false (unsucessful)
func login(handle string, password string) bool {

	return false
}
// Cleanup
func cleanup() {

}

// UI 
func ui() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	//Main UI loop

	quit := false;
	loggedin := false;

	for !quit {
		for !loggedin {
			loggedin,quit = loginWindow();
			if(quit) {
				break
			}
		}
	}
	quitWindow()
}

func loginWindow() (bool, bool) {
	clear()

	write_center(10, "Login:")
	flush()
	time.Sleep(1*time.Second)
	return false, true
}

func quitWindow() {

}

// Termbox functions
func clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func flush() {
	termbox.Flush()
}

func write_center(y int, mess string) {
	x,_ := termbox.Size()
	write_us( ((x/2)-(len(mess)/2)) , y, mess)
}

// Display text on the screen starting at x,y
// Assumes that you are not going to go outside of bounds
func write_us(x int, y int, mess string) {
	for _, c := range mess {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}
// Displays text on the screen starting at x,y and cuts the end off
func write(x int, y int, mess string) {
	sx,_ := termbox.Size()
	if(x+len(mess) > sx) {
		mess = mess[:sx]
	}
	for _, c := range mess {
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

// Display a message in the center of the screen.
func message_us(mess string) {
	_,y := termbox.Size()
	write_center(y/2, mess)
}
