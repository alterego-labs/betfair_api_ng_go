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
	"encoding/json"
	"fmt"
	curl "github.com/andelf/go-curl"
	"log"
)

var loginUrl = "https://identitysso.betfair.com/api/certlogin"

var apiUrl = "https://developer.betfair.com/api.betfair.com/exchange/betting/json-rpc/v1"

type Config struct {
	Username       string
	Password       string
	CertFile       string
	KeyFile        string
	Exchange       string
	Locale         string
	ApplicationKey string
}

type certLoginResult struct {
	LoginStatus  string `json:"loginStatus"`
	SessionToken string `json:"sessionToken"`
}

func CreateSession(c *Config) (certLoginResult, error) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	var result certLoginResult

	writeResultFunc := func(buf []byte, userdata interface{}) bool {
		println("data = >", string(buf))
		err := json.Unmarshal(buf, &result)
		if err != nil {
			log.Fatal("ERROR unmarshal auth data!")
		}
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, writeResultFunc)

	if easy != nil {
		easy.Setopt(curl.OPT_URL, loginUrl)

		easy.Setopt(curl.OPT_HTTPHEADER, []string{"Content-Type: application/x-www-form-urlencoded", "X-Application: " + c.ApplicationKey})
		easy.Setopt(curl.OPT_SSL_VERIFYHOST, 1)
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, true)
		easy.Setopt(curl.OPT_SSLCERT, c.CertFile)
		easy.Setopt(curl.OPT_SSLKEY, c.KeyFile)
		easy.Setopt(curl.OPT_POST, 1)
		easy.Setopt(curl.OPT_POSTFIELDS, "username="+c.Username+"&password="+c.Password)

		code := easy.Perform()

		fmt.Printf("code -> %v\n", code)
	}

	return result, nil
}

func ApiRequest(app_key, session_key, req_data string) string {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	var result string = ""

	writeResultFunc := func(buf []byte, userdata interface{}) bool {
		println("data = >", string(buf))
		result += string(buf)
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, writeResultFunc)

	if easy != nil {
		easy.Setopt(curl.OPT_URL, apiUrl)

		fmt.Printf("request to -> %v\n", apiUrl)

		easy.Setopt(curl.OPT_HTTPHEADER, []string{
			"Content-Type: application/json",
			"Accept: application/json",
			"X-Application: " + app_key,
			"X-Authentication: " + session_key,
		})
		easy.Setopt(curl.OPT_POST, 1)
		easy.Setopt(curl.OPT_POSTFIELDS, req_data)

		code := easy.Perform()

		fmt.Printf("code -> %v\n", code)
	}

	return result
}
