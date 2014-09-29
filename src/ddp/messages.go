package ddp

type m_cConnect struct {
  Msg      string    `json:"msg"`
  Version  string    `json:"version"`
  Support  []string  `json:"support"`
}

type m_sConnected struct {
  Msg      string    `json:"msg"`
  Session  string    `json:"session"`
}

type m_sSub struct {
  Msg      string    `json:"msg"`
  Id       string    `json:"id"`
  Name     string    `json:"name"`
  Params   []string  `json:"params"`
}

type m_sFailed struct {
  Msg      string    `json:"msg"`
  Version  string    `json:"version"`
}

type m_SimpleMessage struct {
  Msg      string    `json:"msg"`
}

type m_sServer struct {
  ServerId  string    `json:"server_id"`
}