package main

type Client struct {
	ID          string
	Name        string
	RedirectUri string
}

type User struct {
	UserID   int
	LoginID  int
	Password string
}

type AuthorizationCode struct {
	Value       string
	UserID      int
	ClientID    int
	Scope       string
	RedirectUri string
	ExpiresAt   string
}

type AccessToken struct {
	Value     string
	UserID    int
	ClientID  int
	Scope     string
	ExpiresAt string
}

type Session struct {
	Client      string
	State       string
	Scope       string
	RedirectUri string
}
