package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
	auth "github.com/sanoyo/authentication"
)

const (
	AUTHORIZATION_CODE_DURATION = 600
	ACCESS_TOKEN_DURATION       = 86400
)

var (
	SUPPORTED_SCOPES = []string{"read", "write"}
	templates        = make(map[string]*template.Template)
	sessionList      = make(map[string]auth.Session)
	authCodeList     = make(map[string]auth.AuthCode)
)

func authorization(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	session := auth.Session{
		Client:      query.Get("client_id"),
		State:       query.Get("state"),
		Scope:       query.Get("scope"),
		RedirectUri: query.Get("redirect_uri"),
	}
	log.Println("session", session)

	// クライアント情報
	clientInfo := auth.Client{
		ID:          "1234",
		Name:        "test",
		RedirectUri: "http://localhost:8080/callback",
	}
	log.Println("clientInfo", clientInfo)

	// paramsのチェック
	if clientInfo.ID != query.Get("client_id") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("client_id is not match"))
		return
	}

	// session id生成
	sessionID := uuid.New().String()
	sessionList[sessionID] = session

	// CookieにセッションIDをセット
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionID,
	}
	http.SetCookie(w, cookie)
	log.Println("cookie", cookie)

	// ログイン&権限認可の画面を戻す
	if err := templates["login"].Execute(w, struct {
		ClientId string
		Scope    string
	}{
		ClientId: session.Client,
		Scope:    session.Scope,
	}); err != nil {
		log.Println(err)
	}
}

// 認可レスポンスを返す
func authCheck(w http.ResponseWriter, req *http.Request) {
	loginUser := req.FormValue("username")
	password := req.FormValue("password")

	// 登録ユーザをハードコード
	user := auth.User{
		UserID:   1111,
		Name:     "hoge",
		Password: "password",
	}

	if loginUser != user.Name || password != user.Password {
		w.Write([]byte("login failed"))
	} else {

		cookie, _ := req.Cookie("session")
		http.SetCookie(w, cookie)
		v, _ := sessionList[cookie.Value]

		authCodeString := uuid.New().String()
		authData := auth.AuthCode{
			User:        loginUser,
			ClientID:    v.Client,
			Scope:       v.Scope,
			RedirectUri: v.RedirectUri,
			ExpiresAt:   time.Now().Unix() + 300,
		}
		// 認可コードを保存
		authCodeList[authCodeString] = authData

		log.Printf("auth code accepet : %s\n", authData)

		location := fmt.Sprintf("%s?code=%s&state=%s", v.RedirectUri, authCodeString, v.State)
		fmt.Println("location", location)

		w.Header().Add("Location", location)
		w.WriteHeader(302)
	}
}

func main() {
	log.Println("start oauth server on localhost:8081...")

	var err error
	templates["login"], err = template.ParseFiles("login.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/authorization", authorization)
	http.HandleFunc("/authcheck", authCheck)
	err = http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
