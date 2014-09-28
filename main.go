package main

import (
  "fmt"
  "bufio"
  "os"
  "ddp"
)

func main() {

  if len(os.Args) != 4 {
    fmt.Println("Usage: go run main.go hostname port path(e.g. 'websocket')")
    os.Exit(1)
  }

  fmt.Println("--> Hit any key to terminate <--")

  ddpClient := ddp.NewDDPClient(os.Args[1], os.Args[2], os.Args[3]) 
  go ddpClient.ListenRead()
  ddpClient.ConnectUsingSaneDefaults()

  reader := bufio.NewReader(os.Stdin)
  text, _ := reader.ReadString('\n')
  fmt.Println(text)
}