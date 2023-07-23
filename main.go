package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("[!] Command needs 3 arguments:")
		fmt.Println("./gha-private-auth $PEM $APP_ID $INSTALL_ID")
		os.Exit(1)
	}

	pemStr := os.Args[1]
	appId := os.Args[2]
	installId := os.Args[3]

	key, err := parsePem(pemStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	id, err := strconv.Atoi(appId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jwt, err := createJwt(key, id)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	token, err := getToken(jwt, installId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(token)
}

func parsePem(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM format, failed decode")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil

}

func createJwt(key *rsa.PrivateKey, appId int) (string, error) {

	t := time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"iat": t,
			"exp": t + 600,
			"iss": appId,
		})

	jwt, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func getToken(jwt, installId string) (string, error) {

	type tokenResponse struct {
		Token string `json:"token"`
	}

	var t tokenResponse

	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installId)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		failed := true
		for i := 1; i < 5; i++ {
			resp, err = client.Do(req)
			if err == nil {
				failed = false
				break
			} else {
				time.Sleep(time.Duration(2 * time.Second))
			}
		}
		if failed {
			return "", err
		}
	}

	defer resp.Body.Close()

	if 200 <= resp.StatusCode && resp.StatusCode <= 299 {
		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			return "", err
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("invalid response status %s\nbody:%s", fmt.Sprint(resp.StatusCode), string(body))
	}

	return t.Token, nil
}
