package ddp

import (
  "log"
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

func (ddpClient *DDPClient) ListenRead() {
  log.Println("Listening read from server")
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
        log.Println("ServerId: " + data.ServerId)
      
      case strings.Contains(msg, "connected"):
        ddpClient.connected = true
        var data m_sConnected
        if err := json.Unmarshal([]byte(msg), &data); err != nil {
          panic(err)
        }
        log.Println("Connected to ddp server with session: " + data.Session)
      
      case strings.Contains(msg, "failed"):
        var data m_sFailed
        if err := json.Unmarshal([]byte(msg), &data); err != nil {
          panic(err)
        }
        log.Println(data.Msg + ", server requires version: " + data.Version)
      // <- Establishing a DDP Connection
      
      // Heartbeat ->
      case strings.Contains(msg, "ping"):
        log.Println("Received ping, sending pong")
        ddpClient.sendSimpleMessage("pong")
      // <- Heartbeat
      
      // Managing Data ->
      case strings.Contains(msg, "nosub"):
        log.Println("received: " + "'nosub'");
      case strings.Contains(msg, "added"):
        log.Println("received: " + "'added'");
      case strings.Contains(msg, "changed"):
        log.Println("received: " + "'changed'");
      case strings.Contains(msg, "removed"):
        log.Println("received: " + "'removed'");
      case strings.Contains(msg, "ready"):
        log.Println("received: " + "'ready'");
      case strings.Contains(msg, "addedBefore"):
        log.Println("received: " + "'addedBefore'");
      case strings.Contains(msg, "movedBefore"):
        log.Println("received: " + "'movedBefore'");
      // <- Managing Data
      
      // Remote Procedure Calls ->
      case strings.Contains(msg, "result"):
        log.Println("received: " + "'result'");
      case strings.Contains(msg, "updated"):
        log.Println("received: " + "'updated'");
      // <- Remote Procedure Calls

    }
  }
}