// Non-interactive authentication in betfair API-NG
// Copyright (C) 2014  Sergey Gernyak

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package betfair_api_ng_go

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	// "net"
	"net/http"
	"strings"
	// "time"
)

var loginUrl = "https://identitysso.betfair.com/api/certlogin"

type Config struct {
	Username string
	Password string
	CertFile string
	KeyFile  string
	Exchange string
	Locale   string
}

type certLoginResult struct {
	LoginStatus  string `json:"loginStatus"`
	SessionToken string `json:"sessionToken"`
}

func CreateSession(c *Config) (string, error) {
	body := strings.NewReader("username=" + c.Username + "&password=" + c.Password)

	data, err := doRequest("certLogin", "", body, c)
	// log.Fatal("DATA " + data)
	var result certLoginResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal("ERROR unmarshal!")
		return "", err
	}
	if result.LoginStatus != "SUCCESS" {
		return "", errors.New(result.LoginStatus)
	}

	return result.SessionToken, nil
}

func doRequest(key, method string, body *strings.Reader, c *Config) ([]byte, error) {

	req, err := http.NewRequest("POST", loginUrl, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	if key == "certLogin" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// In non-interactive login, X-Application is not validated
		req.Header.Set("X-Application", "5kaGlzvjvo8HeaNo")
	}

	res, err := httpClient(c).Do(req)
	if err != nil {
		fmt.Printf("%s\n", err)
		log.Fatal("ERROR do!")
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Fatal("ERROR with status code!")
		return nil, errors.New(res.Status)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("ERROR readall!")
		return nil, err
	}

	return data, nil
}

func httpClient(c *Config) (httpClient *http.Client) {
	cert, _ := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)

	ssl := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	ssl.Rand = rand.Reader
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: ssl,
		},
	}
}
