package sample

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type Json struct {
	Data string `json:"data"`
}

var (
	authClient *auth.Client
)

func init() {
	var err error

	config := &firebase.Config{
		ProjectID: os.Getenv("PROJECT_ID"),
	}
	firebaseApp, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err = firebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("firebaseApp.Auth: %v", err)
	}
}

func OnCall(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// ユーザ認証
	// Authorization: Bearer [token]  (RFC 6750)
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]
	t, err := authClient.VerifyIDToken(ctx, token)
	if err != nil {
		log.Fatalf("verify token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// ユーザID表示
	userId := t.Claims["user_id"]
	log.Printf("user_id = %s", userId)

	// リクエスト読み込み
	body, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		log.Fatalf("read body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// JSONパース
	var requestJson Json
	err = json.Unmarshal(body, &requestJson)
	if err != nil {
		log.Fatalf("json unmarshal: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// リクエストデータ表示
	fmt.Println(requestJson.Data)

	// JSONエンコード
	responseJson, err := json.Marshal(map[string]string{
		"data": "pong",
	})
	log.Printf("data: %s", responseJson)
	if err != nil {
		log.Fatalf("marshal json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// レスポンス出力
	fmt.Fprint(w, string(responseJson))
}
