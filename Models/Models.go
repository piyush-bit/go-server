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

type ForgetPassword struct {
	Email string `json:"email"`
	Token string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}


type GoogleUser struct {
    Iss            string `json:"iss"`
    Azp            string `json:"azp"`
    Aud            string `json:"aud"`
    Sub            string `json:"sub"`
    Email          string `json:"email"`
    EmailVerified  string `json:"email_verified"`
    Name           string `json:"name"`
    Picture        string `json:"picture"`
    GivenName      string `json:"given_name"`
    FamilyName     string `json:"family_name"`
    Locale         string `json:"locale"`
}