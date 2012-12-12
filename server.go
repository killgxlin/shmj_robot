package main
import (
	"fmt"
	"os"
	"net"
	"bufio"
	"strings"
)

func main(){
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	input := make(chan string, 100)
	output := make(chan string, 100)
	go func (){
		bufReader := bufio.NewReader(conn)
		for {
			readed, err := bufReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			
			input<-strings.TrimSpace(readed)
		}
	} ()
	go func() {
		for {
			msg := <-output
			_, err := conn.Write([]byte(msg + "\r\n")) 
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	} ()

	for {
		select {
			case msg := <-input:
				output<-msg;
		}
	}
}
