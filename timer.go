package main
import (
	"time"
	"fmt"
)

func main() {

	timeout := make(chan bool, 3)
	go func(){
		for {
			time.Sleep(1 * time.Second)
			timeout<-true
		}
	}()
	for {
		select {
			case <-timeout:
				fmt.Println("hello world")
		}
	}
}
