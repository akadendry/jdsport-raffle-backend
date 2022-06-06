package models

type ResponseLogin struct {
    ID    string    `json:"id"`
	Message string    `json:"message"`
	AppName string    `json:"appName"`
	Datas Tokend `json:"data"`
}

type Tokend struct {
	Token     string `json:"token"`
	Expired        int    `json:"expired"`
}
