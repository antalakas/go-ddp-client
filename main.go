package ddp

import (
	"code.google.com/p/go.net/websocket"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"time"
)

type DDPClient struct {
	host            string
	port            int
	path            string
	connected       bool
	sessionId       string
	pingSeconds     time.Duration
	pongWaitSeconds time.Duration
	ws              *websocket.Conn
	sub_id          int
	next_id         int
	readyChan       chan bool
	loginUser       bool
	enableLogin     bool
	callbacks       map[string]interface{}
}

// Create new chat client.
func NewDDPClient(host string, port string,
	path string, enableLogin bool) *DDPClient {
	var ddpClient DDPClient
	ddpClient.connectSocket(host, port, path)

	if ddpClient.ws == nil {
		panic("There is no websocket")
	}

	ddpClient.connected = false
	ddpClient.sessionId = ""
	ddpClient.pingSeconds = time.Second * 25
	ddpClient.pongWaitSeconds = time.Second * 5
	ddpClient.sub_id = 0
	ddpClient.next_id = 0
	ddpClient.loginUser = true
	ddpClient.enableLogin = enableLogin
	ddpClient.callbacks = make(map[string]interface{})
	return &ddpClient
}

func (ddpClient *DDPClient) NextId() int {
	ddpClient.next_id = ddpClient.next_id + 1
	return ddpClient.next_id
}

func (ddpClient *DDPClient) connectSocket(host string, port string, path string) {
	conn, err := websocket.Dial("ws://"+host+":"+
		string(port)+"/"+
		path, "", "http://"+host)

	checkError("connectSocket", err)
	ddpClient.ws = conn
}

func (ddpClient *DDPClient) ConnectUsingSaneDefaults(readyChan chan bool) {
	var msg = &m_cConnect{"connect", "pre1", []string{"1", "pre2", "pre1"}}
	ddpClient.readyChan = readyChan
	websocket.JSON.Send(ddpClient.ws, msg)
}

func (ddpClient *DDPClient) Connect(version string, support []string) {
	var msg = &m_cConnect{"connect", version, support}
	websocket.JSON.Send(ddpClient.ws, msg)
}

func (ddpClient *DDPClient) LoginUser() {

	if !ddpClient.loginUser {
		return
	}

	//printMessage("MY_USERNAME:" + os.Getenv("MY_USERNAME"))
	//printMessage("MY_PASSWORD:" + os.Getenv("MY_PASSWORD"))

	hash := sha256.New()
	hash.Write([]byte(os.Getenv("MY_PASSWORD")))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	//printMessage(mdStr)

	var userMsg = m_Username{os.Getenv("MY_USERNAME")}
	var passMsg = m_Password{mdStr, "sha-256"}
	var credMsg = m_UserCredentials{userMsg, passMsg}
	var loginMsg = &m_cUserLogin{"method", "login",
		[]m_UserCredentials{credMsg}, "2"}

	//loginMsgStr, _ := json.Marshal(loginMsg)
	//printMessage("Login message: ")
	//fmt.Printf("\nMarshalled data: %s\n", loginMsgStr)

	websocket.JSON.Send(ddpClient.ws, loginMsg)
}

func (ddpClient *DDPClient) LoginEmail() {

	if ddpClient.loginUser {
		return
	}

	hash := sha256.New()
	hash.Write([]byte(os.Getenv("MY_PASSWORD")))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)

	var emailMsg = m_Email{os.Getenv("MY_EMAIL")}
	var passMsg = m_Password{mdStr, "sha-256"}
	var credMsg = m_EmailCredentials{emailMsg, passMsg}
	var loginMsg = &m_cEmailLogin{"method", "login",
		[]m_EmailCredentials{credMsg}, "2"}

	websocket.JSON.Send(ddpClient.ws, loginMsg)
}

func (ddpClient *DDPClient) Logout() {
	var logoutMsg = &m_cLogout{"method", "logout", []string{}, "2"}
	websocket.JSON.Send(ddpClient.ws, logoutMsg)
}

func (ddpClient *DDPClient) Subscribe(name string) {
	var msg = &m_sSub{"sub", string(ddpClient.sub_id), name, []string{}}
	ddpClient.sub_id = ddpClient.sub_id + 1
	websocket.JSON.Send(ddpClient.ws, msg)
}

func (ddpClient *DDPClient) Unsubscribe() {

}

func (ddpClient *DDPClient) sendSimpleMessage(msg string) {
	websocket.JSON.Send(ddpClient.ws, &m_SimpleMessage{msg})
}

func (ddpClient *DDPClient) CallMethod(method string, params []interface{}, callback interface{}) {
	id := string(ddpClient.NextId())
	ddpClient.callbacks[id] = callback
	websocket.JSON.Send(ddpClient.ws, &m_RPC{"method", method, params, id})
}

func (ddpClient *DDPClient) pingJob(pongChan chan bool) {

	pingTicker := time.NewTicker(ddpClient.pingSeconds).C

	for {
		select {
		case <-pingTicker:
			printMessage("sending ping")
			ddpClient.sendSimpleMessage("ping")
			pongTimer := time.NewTimer(ddpClient.pongWaitSeconds)
			pongTimerChan := pongTimer.C
			select {
			case <-pongTimerChan:
				ClientExit("I've waited enough for pong, exiting...")
			case <-pongChan:
				pongTimer.Stop()
				printMessage("received: " + "'pong'")
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

		if msg == "" {
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

			if ddpClient.enableLogin {
				ddpClient.LoginUser()
			}

			ddpClient.Subscribe("meteor.loginServiceConfiguration")
			ddpClient.Subscribe("meteor_autoupdate_clientVersions")
			ddpClient.readyChan <- true

		case strings.Contains(msg, "failed"):
			var data m_sFailed
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				panic(err)
			}
			printMessage(data.Msg + ", server requires version: " + data.Version)
		// <- Establishing a DDP Connection

		// Heartbeat ->
		case strings.Contains(msg, "ping"):
			printMessage("received: " + "'ping'" + ", sending 'pong'")
			ddpClient.sendSimpleMessage("pong")
		case strings.Contains(msg, "pong"):
			pongChan <- true
			//printMessage("received: " + "'pong'");
		// <- Heartbeat

		// Managing Data ->
		case strings.Contains(msg, "nosub"):
			printMessage("received: " + "'nosub'")
		case strings.Contains(msg, "added"):
			printMessage("received: " + "'added'")
			printMessage("Message received: " + msg)
		case strings.Contains(msg, "changed"):
			printMessage("received: " + "'changed'")
		case strings.Contains(msg, "removed"):
			printMessage("received: " + "'removed'")
		case strings.Contains(msg, "ready"):
			printMessage("received: " + "'ready'")
			printMessage("Message received: " + msg)
		case strings.Contains(msg, "addedBefore"):
			printMessage("received: " + "'addedBefore'")
		case strings.Contains(msg, "movedBefore"):
			printMessage("received: " + "'movedBefore'")
		// <- Managing Data

		// Remote Procedure Calls ->
		case strings.Contains(msg, "result"):
			printMessage("received: " + "'result'")
			var data m_RPCResult
			if err := json.Unmarshal([]byte(msg), &data); err != nil {
				panic(err)
			}
			ddpClient.callbacks[data.Id].(func(interface{}, interface{}))(data.Error, data.Result)
			delete(ddpClient.callbacks, data.Id)
		case strings.Contains(msg, "updated"):
			printMessage("received: " + "'updated'")
			// <- Remote Procedure Calls

		}
	}
}
