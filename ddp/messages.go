package ddp

type m_cConnect struct {
	Msg     string   `json:"msg"`
	Version string   `json:"version"`
	Support []string `json:"support"`
}

type m_sConnected struct {
	Msg     string `json:"msg"`
	Session string `json:"session"`
}

type m_sSub struct {
	Msg    string   `json:"msg"`
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Params []string `json:"params"`
}

type m_sFailed struct {
	Msg     string `json:"msg"`
	Version string `json:"version"`
}

type m_SimpleMessage struct {
	Msg string `json:"msg"`
}

type m_sServer struct {
	ServerId string `json:"server_id"`
}

type m_Username struct {
	Username string `json:"username"`
}

type m_Email struct {
	Email string `json:"email"`
}

type m_Password struct {
	Digest    string `json:"digest"`
	Algorithm string `json:"algorithm"`
}

type m_UserCredentials struct {
	User     m_Username `json:"user"`
	Password m_Password `json:"password"`
}

type m_EmailCredentials struct {
	User     m_Email    `json:"user"`
	Password m_Password `json:"password"`
}

type m_cUserLogin struct {
	Msg    string              `json:"msg"`
	Method string              `json:"method"`
	Params []m_UserCredentials `json:"params"`
	Id     string              `json:"id"`
}

type m_cEmailLogin struct {
	Msg    string               `json:"msg"`
	Method string               `json:"method"`
	Params []m_EmailCredentials `json:"params"`
	Id     string               `json:"id"`
}

type m_cLogout struct {
	Msg    string   `json:"msg"`
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     string   `json:"id"`
}
