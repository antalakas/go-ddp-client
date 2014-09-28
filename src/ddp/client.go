package ddp

import (
  "time"
  "strings"
  "encoding/json"
  "code.google.com/p/go.net/websocket"
)

type DDPClient struct {
  host string
  port int
  path string
  connected bool
  sessionId string
  pingSeconds time.Duration
  pongWaitSeconds time.Duration
  ws *websocket.Conn
}

// Create new chat client.
func NewDDPClient(host string, port string, path string) *DDPClient {
  var ddpClient DDPClient
  ddpClient.connectSocket(host, port, path)
  
  if ddpClient.ws == nil {
    panic("There is no websocket")
  }
  
  ddpClient.connected = false
  ddpClient.sessionId = ""
  ddpClient.pingSeconds = time.Second * 25
  ddpClient.pongWaitSeconds = time.Second * 5
  
  return &ddpClient
}

func (ddpClient *DDPClient) connectSocket(host string, port string, path string) {             
  conn, err := websocket.Dial("ws://" + host + ":" + 
                              string(port) + "/" + 
                              path, "", "http://" + host)
  
  checkError("connectSocket", err)
  ddpClient.ws = conn
}

func (ddpClient *DDPClient) ConnectUsingSaneDefaults() {
  var msg = &m_cConnect{"connect", "1", []string{"1","pre2","pre1"}}
  //connectMsg, _ := json.Marshal(msg)
  //fmt.Println("Connect message: ")
  //fmt.Printf("\nMarshalled data: %s\n", connectMsg)
  websocket.JSON.Send(ddpClient.ws, msg)
}

func (ddpClient *DDPClient) Connect(version string, support []string) {
  var msg = &m_cConnect{"connect", version, support}
  websocket.JSON.Send(ddpClient.ws, msg)
}

func (ddpClient *DDPClient) Login() {

}

func (ddpClient *DDPClient) Logout() {

}

func (ddpClient *DDPClient) Subscribe() {

}

func (ddpClient *DDPClient) Unsubscribe() {

}

func (ddpClient *DDPClient) sendSimpleMessage(msg string) {
  websocket.JSON.Send(ddpClient.ws, &m_SimpleMessage{msg})
}

func (ddpClient *DDPClient) pingJob(pongChan chan bool) {

  pingTicker := time.NewTicker(ddpClient.pingSeconds).C
  
  for {
    select {
      case <- pingTicker:
        ddpClient.sendSimpleMessage("ping")
        pongChan := time.NewTimer(ddpClient.pongWaitSeconds).C
        select {
        case <- pongChan:
            clientExit("I've waited enough for pong, exiting...")
        case <- pongChan:
            printMessage("received: " + "'pong'");
        }
    }
  }
  
}

func (ddpClient *DDPClient) ListenRead() {
  
  pongChan := make(chan bool)
  
  printMessage("Listening read from server")
  for {
    var msg string
    websocket.Message.Receive(ddpClient.ws, &msg)

    if (msg == "") {
      continue
    }

    //fmt.Println("Message received: " + msg)
    
    switch {
      
      // Establishing a DDP Connection ->
      case strings.Contains(msg, "server_id"):
        var data m_sServer
        if err := json.Unmarshal([]byte(msg), &data); err != nil {
          panic(err)
        }
        printMessage("ServerId: " + data.ServerId)
      
      case strings.Contains(msg, "connected"):
        ddpClient.connected = true
        var data m_sConnected
        if err := json.Unmarshal([]byte(msg), &data); err != nil {
          panic(err)
        }
        printMessage("Connected to ddp server with session: " + data.Session)
        
        go ddpClient.pingJob(pongChan)
      
      case strings.Contains(msg, "failed"):
        var data m_sFailed
        if err := json.Unmarshal([]byte(msg), &data); err != nil {
          panic(err)
        }
        printMessage(data.Msg + ", server requires version: " + data.Version)
      // <- Establishing a DDP Connection
      
      // Heartbeat ->
      case strings.Contains(msg, "ping"):
        printMessage("received: " + "'ping'" + ", sending 'pong'");
        ddpClient.sendSimpleMessage("pong")
      case strings.Contains(msg, "pong"):
        pongChan <- true
        //printMessage("received: " + "'pong'");
      // <- Heartbeat
      
      // Managing Data ->
      case strings.Contains(msg, "nosub"):
        printMessage("received: " + "'nosub'");
      case strings.Contains(msg, "added"):
        printMessage("received: " + "'added'");
      case strings.Contains(msg, "changed"):
        printMessage("received: " + "'changed'");
      case strings.Contains(msg, "removed"):
        printMessage("received: " + "'removed'");
      case strings.Contains(msg, "ready"):
        printMessage("received: " + "'ready'");
      case strings.Contains(msg, "addedBefore"):
        printMessage("received: " + "'addedBefore'");
      case strings.Contains(msg, "movedBefore"):
        printMessage("received: " + "'movedBefore'");
      // <- Managing Data
      
      // Remote Procedure Calls ->
      case strings.Contains(msg, "result"):
        printMessage("received: " + "'result'");
      case strings.Contains(msg, "updated"):
        printMessage("received: " + "'updated'");
      // <- Remote Procedure Calls

    }
  }
}