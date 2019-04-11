package user

type User struct {
	Name       string `json:"name"`
	Account_ID string `json:"accountid"`
	API_KEY    string `json:"apikey"`
	API_SECRET string `json:"apisecret"`
}
