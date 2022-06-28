package model

type Favorite struct {
	Id       string
	Favorite string
}

type ShortMember struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Member struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Favorites Favorites `json:"favorites"`
}

type Favorites []string

type Message struct {
	Result string `json:"result"`
}

var ResultOK Message

func init() {
	ResultOK.Result = "OK"
}
