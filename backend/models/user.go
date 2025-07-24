package models

type User struct {
	ID       string    `json:"id"`
	Username string 	 `json:"username"`
	Email    string 	 `json:"email"`
	Password []byte 	 `json:"-"`
	Name 		 string 	 `json:"name"`
	Surname  string 	 `json:"surname"`
}
