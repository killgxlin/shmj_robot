package main
import (
	"fmt"
	"time"
	"bufio"
	"os"
	"strings"
	"strconv"
)

type IO struct {
	input chan string
	output chan string
}

func newIO(a chan string, b chan string) *IO {
	if a==nil {
		a = make(chan string, 100)
	}
	if b==nil {
		b = make(chan string, 100)
	}
	return &IO{input:a,	output:b}
}

/*	
	start 1 100
	shut 10
	logon 10
	combat 10
	chat 10
*/
func start_console() *IO {
	console := newIO(nil, nil)
	go func(input chan string) {
		for {
			select {
				case msg:=<-input:
					fmt.Println(msg)
			}
		}
	}(console.input)

	go func(output chan string) {
		consoleBuffer := bufio.NewReader(os.Stdin)
		for {
			line, err := consoleBuffer.ReadString('\n')
			checkError(err)

			if command := strings.Trim(line, "\r\n "); command != "" {
				output<-command
			}
		}
	}(console.output)

	return console
}

func checkError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func report_stat(output chan string, playeridstr string, srcStat string, dstStat string) {
	output<-playeridstr + " stat " + srcStat + " to " + dstStat
}

func Connect(io *IO) {
	conn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Println(err)
		return
	}
	readwriter := bufio.NewReadWriter(conn, conn)
	go func(){
		for {
			msg, err := readwriter.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			io.input<-strings.TrimSpace(msg)
		}
	}()
	go func(){
		for {
			select {
				case msg := <-conn.output:
					len, err := readwriter.WriteString(msg)
					if err != nil {
						fmt.Println(err)
						return
					}
					err = readwriter.Flush()
					if err != nil {
						fmt.Println(err)
						return
					}
			}
		}
	}()
}
/*
	offline stat
		shut
		logon
	online stat
		shut
		logoff
		chat
		combat
	
*/
func start_robot(output chan string, playerid uint64) *Robot {
	ioer := newIO(nil, output)
	playeridstr := strconv.FormatUint(playerid, 10)

	go func() {
		nextStat := "offline"
		chat := false
		combat := false
		timer := make(chan bool)
		go func() {
			time.Sleep(1 * time.Second)
			timer<-true
		}()
		stat := "null"
		conn := newIO(nil, nil)
		FOR: for {
			switch stat {
				case "offline":
					select {
						case cmd := <-ioer.input:
							args := strings.Split(strings.TrimSpace(cmd))
							switch args[0] {
								case "shut":
									nextStat = "null"
								case "logon":
									Connect(conn)
									conn.output<-CSLogon
							}
						case msg := <-conn.input:
							args := strings.Split(strings.TrimSpace(msg))
							switch args[0] {
								case SCLogon:
									nextStat = "online"
							}
						case <-timer:
					}
				case "online":
					select {
						case cmd := <-ioer.input:
							args := strings.Split(strings.Trim(cmd, "\r\n\t "), " ")
							switch args[0] {
								case "shut":
									nextStat = "null"
								case "logoff":
									conn.Close
									nextStat = "offline"
								case "chat":
									chat = !chat 
								case "combat":
									combat = !combat
							}
						case msg := <-conn.input:
							args := strings.Split(strings.TrimSpace(msg))
							switch args[0] {
							}
						case <-timer:
							if chat {
							}
							if combat {
							}
					}
				case "null":
					switch nextStat {
						case "null":
							break FOR
						default:
							
					}
			}
			if nextStat != stat {
				report_stat(ioer.output, playeridstr, stat, nextStat)
				stat = nextStat
			}
		}

	}()

	return &Robot{IO:ioer, stat:"null"}
}

type Robot struct {
	*IO
	stat string
}
type Robots map[uint64]*Robot
type Operator func(playerid uint64, robot *Robot)bool
func forEach(stat string, robots Robots, op Operator) {
	for k,v := range robots {
		if v != nil && ( v.stat == stat || stat == "" ){
			if op(k, v) {
				break
			}
		}
	}
}
			

func main() {
	console := start_console()
	
	report := make(chan string, 100)
	robots := make(Robots)

	FOR:for {
		select {
			case msg := <-console.output:
				args := strings.Split(strings.Trim(msg, "\r\n\t "), " ")
				switch args[0] {
					case "start":
						start, err := strconv.ParseInt(args[1], 10, 32)
						checkError(err)
						arg2, err := strconv.ParseInt(args[2], 10, 32)
						checkError(err)
						for i:=start; i<=arg2; i++ {
							playerid := uint64(i)
							if _, found := robots[playerid]; !found {
								robots[playerid] = start_robot(report, playerid)
							}
						}
					case "list":
						list := "( "
						forEach(args[1], robots, func(playerid uint64, robot *Robot) bool {
							list += (strconv.FormatUint(playerid, 10) + " ")
							return false
						})
						list += " )"
						console.input<-list
					case "shut":
						num, err := strconv.ParseInt(args[1], 10, 32)
						checkError(err)
						forEach("", robots, func(playerid uint64, robot *Robot) bool {
							robot.input<-"shut"
							num--
							return num <= 0
						})
					case "quit":
						break FOR
					case "logon":
						num, err := strconv.ParseInt(args[1], 10, 32)
						checkError(err)
						forEach("", robots, func(playerid uint64, robot *Robot) bool {
							robot.input<-"logon"
							num--
							return num <= 0
						})
					default:
						console.input<-"invalid command"
				}
			case msg:=<-report:
				console.input<-msg
				args := strings.Split(msg, " ")
				playerid, err := strconv.ParseUint(args[0], 10, 64)
				checkError(err)
				switch args[1] {
					case "shut":
					case "stat":
						switch args[4] {
							case "null":
								delete(robots, playerid)
							default:
								robots[playerid].stat = args[4]
						}	
						
				}
		}
	}
}

