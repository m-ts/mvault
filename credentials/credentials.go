package credentials

/*Credentials stores password and salt for encryption*/
type Credentials struct {
	Password string `json:"password"`
	Salt     string `json:"salt"`
}
