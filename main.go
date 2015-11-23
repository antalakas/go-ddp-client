package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"./ddp"
)

func application(ddpClient *ddp.DDPClient) {
	fmt.Println("")
	fmt.Println("-----------------")
	fmt.Println("Hey, i am running")

	ddpClient.Logout()
}

func entryPoint(ddpClient *ddp.DDPClient, readyChan chan bool) {

	// Wait a second, to init
	time.Sleep(time.Second)

	initTimer := time.NewTimer(time.Second * 10)
	initTimerChan := initTimer.C

	for {
		select {
		case <-initTimerChan:
			ddp.ClientExit("I've waited enough for initialization, exiting...")
		case <-readyChan:
			initTimer.Stop()
			fmt.Println("Done init")
			go application(ddpClient)
		}
	}

}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go hostname port path(e.g. 'websocket')")
		os.Exit(1)
	}

	fmt.Println("--> Hit any key to terminate <--")

	// last argument is used to enable login
	ddpClient := ddp.NewDDPClient(os.Args[1], os.Args[2], os.Args[3], true)
	go ddpClient.ListenRead()

	readyChan := make(chan bool)
	go entryPoint(ddpClient, readyChan)

	ddpClient.ConnectUsingSaneDefaults(readyChan)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}
