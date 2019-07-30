package auth

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/router"
)

// ReqGetBasic Struct
type ReqGetBasic struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Basic Function as Midleware for Basic Authorization
func Basic(next http.Handler) http.Handler {
	// Return Next HTTP Handler Function, If Authorization is Valid
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse HTTP Header Authorization
		authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		// Check HTTP Header Authorization Section
		// Authorization Section Length Should Be 2
		// The First Authorization Section Should Be "Basic"
		if len(authHeader) != 2 || authHeader[0] != "Basic" {
			hlp.LogPrintln(hlp.LogLevelWarn, "http-access", "unauthorized method "+r.Method+" at URI "+r.RequestURI)
			router.ResponseAuthenticate(w)
			return
		}

		// The Second Authorization Section Should Be The Credentials Payload
		// But We Should Decode it First From Base64 Encoding
		authPayload, err := base64.StdEncoding.DecodeString(authHeader[1])
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

		// Split Decoded Authorization Payload Into Username and Password Credentials
		authCredentials := strings.SplitN(string(authPayload), ":", 2)

		// Check Credentials Section
		// It Should Have 2 Section, Username and Password
		if len(authCredentials) != 2 {
			hlp.LogPrintln(hlp.LogLevelWarn, "http-access", "unauthorized method "+r.Method+" at URI "+r.RequestURI)
			router.ResponseBadRequest(w, "")
			return
		}

		// Make Credentials to JSON Format
		jsonCredentials := `{"username": "` + authCredentials[0] + `", "password": "` + authCredentials[1] + `"}`

		// Rewrite Body Content With Credentials in JSON Format
		r.Body = ioutil.NopCloser(strings.NewReader(jsonCredentials))

		// Call Next Handler Function With Current Request
		next.ServeHTTP(w, r)
	})
}
