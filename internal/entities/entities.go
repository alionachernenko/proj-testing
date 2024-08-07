package entities

type Task struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Done  bool   `json:"done"`
}

type User struct {
	Username string
	Password string
}
