package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
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
	ID           int    `json:"id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	AppId        int    `json:"app_id"`
}

type Session struct {
	ID       int    `json:"id"`
	RefreshToken    string `json:"token"`
	UserId   int    `json:"user_id"`
	AppId    int    `json:"app_id"`
}