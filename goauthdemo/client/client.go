package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/xuperdata/xuperdid/demo/jwtutil"
)

const (
	authServerURL = "http://localhost:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/authorize",
			TokenURL: authServerURL + "/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		u := config.AuthCodeURL("xyz")
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc(
		"/oauth2",
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			state := r.Form.Get("state")
			if state != "xyz" {
				http.Error(w, "State invalid", http.StatusBadRequest)
				return
			}
			code := r.Form.Get("code")
			if code == "" {
				http.Error(w, "Code not found", http.StatusBadRequest)
				return
			}
			println("duanbing: " + code)
			token, err := config.Exchange(context.Background(), code)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			globalToken = token

			e := json.NewEncoder(w)
			e.SetIndent("", "  ")
			e.Encode(token)
		})

	http.HandleFunc(
		"/refresh",
		func(w http.ResponseWriter, r *http.Request) {
			if globalToken == nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			globalToken.Expiry = time.Now()
			token, err := config.TokenSource(context.Background(), globalToken).Token()
			if err != nil {
				http.Error(w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			globalToken = token
			e := json.NewEncoder(w)
			e.SetIndent("", "  ")
			e.Encode(token)
		})
	http.HandleFunc(
		"/verify",
		func(w http.ResponseWriter, r *http.Request) {
			if globalToken == nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			//signedKey := []byte(jwtutil.SIGNED_KEY)
			signedKey := jwtutil.GetPublicKey()
			token, err := jwt.Parse(globalToken.AccessToken, func(token *jwt.Token) (interface{}, error) {
				//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return signedKey, nil
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			parts := strings.Split(globalToken.AccessToken, ".")
			verifyErr := token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, signedKey)
			if verifyErr != nil {
				fmt.Println(verifyErr)
				e := json.NewEncoder(w)
				e.SetIndent("", "  ")
				e.Encode(verifyErr.Error())
				return
			}
			if claims, ok := token.Claims.(jwt.Claims); ok && token.Valid {
				fmt.Printf("claims: %#v", claims)
				e := json.NewEncoder(w)
				e.SetIndent("", "  ")
				e.Encode(claims)
			} else {
				fmt.Println(err)
			}
		})

	http.HandleFunc(
		"/try",
		func(w http.ResponseWriter, r *http.Request) {
			if globalToken == nil {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}

			resp, err := http.Get(fmt.Sprintf(
				"%s/test?access_token=%s", authServerURL,
				globalToken.AccessToken))
			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusBadRequest)
				return
			}
			defer resp.Body.Close()

			io.Copy(w, resp.Body)
		})

	http.HandleFunc(
		"/pwd",
		func(w http.ResponseWriter, r *http.Request) {
			token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			globalToken = token
			e := json.NewEncoder(w)
			e.SetIndent("", "  ")
			e.Encode(token)
		})

	http.HandleFunc(
		"/client",
		func(w http.ResponseWriter, r *http.Request) {
			cfg := clientcredentials.Config{
				ClientID:     config.ClientID,
				ClientSecret: config.ClientSecret,
				TokenURL:     config.Endpoint.TokenURL,
			}

			token, err := cfg.Token(context.Background())
			if err != nil {
				http.Error(
					w, err.Error(),
					http.StatusInternalServerError)
				return
			}

			e := json.NewEncoder(w)
			e.SetIndent("", "  ")
			e.Encode(token)
		})

	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
