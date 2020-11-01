package credentials

/*Credentials stores password and salt for encryption*/
type Credentials struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}

/*RequestCredentials stores password and salt for encryption*/
type RequestCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
