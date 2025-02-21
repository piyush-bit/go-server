package models

type User struct {
	ID       int    `json:"id"`
	Name	 string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type App struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CallbackUrl string `json:"callback_url"`
	UserId      int    `json:"user_id"`
}

type Token struct {
	ID        int    `json:"id"`
	Token     string `json:"token"`
	AppId     int    `json:"app_id"`
}


