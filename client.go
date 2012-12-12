package main
import (
	"net"
	"fmt"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func(){
		for {
			msg := make([]byte, 1000)
			len, err := conn.Read(msg)
			if err != nil {
				fmt.Println(err)
				break
			}
			_, err = os.Stdout.Write(msg[0:len])
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}()		

	for {
		msg := make([]byte, 1000)
		len, err := os.Stdin.Read(msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		_, err = conn.Write(msg[0:len])
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
