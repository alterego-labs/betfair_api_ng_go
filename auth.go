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
	curl "github.com/juliuxu/go-curl"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var loginUrl = "https://identitysso.betfair.com/api/certlogin"

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

func CreateSession(c *Config) (string, error) {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	fooTest := func(buf []byte, userdata interface{}) bool {
		println("size=>", len(buf))
		println("DEBUG(in callback)", buf, userdata)
		println("data = >", string(buf))
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if easy != nil {
		easy.Setopt(curl.OPT_URL, loginUrl)

		easy.Setopt(curl.OPT_HTTPHEADER, []string{"Content-Type: application/x-www-form-urlencoded", "X-Application: " + c.ApplicationKey})
		easy.Setopt(curl.OPT_HEADER, 1)

		easy.Setopt(curl.OPT_VERBOSE, true)

		easy.Setopt(curl.OPT_SSL_VERIFYHOST, 1)
		easy.Setopt(curl.OPT_SSL_VERIFYPEER, true)
		easy.Setopt(curl.OPT_SSLCERT, c.CertFile)
		easy.Setopt(curl.OPT_SSLKEY, c.KeyFile)
		easy.Setopt(curl.OPT_POST, 1)
		easy.Setopt(curl.OPT_POSTFIELDS, "username="+c.Username+"&password="+c.Password)

		code := easy.Perform()

		fmt.Printf("code -> %v\n", code)
	}

}
