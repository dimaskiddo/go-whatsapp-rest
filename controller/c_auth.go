package controller

import (
	"encoding/json"
	"net/http"

	svc "github.com/dimaskiddo/whatsapp-go-rest/service"
)

// GetAuth Function to Get Authorization Token
func GetAuth(w http.ResponseWriter, r *http.Request) {
	var creds svc.BasicCredentials
	_ = json.NewDecoder(r.Body).Decode(&creds)

	if len(creds.Username) == 0 || len(creds.Password) == 0 {
		svc.ResponseBadRequest(w, "Invalid authorization")
		return
	}

	if creds.Password == "83e4060e-78e1-4fe5-9977-aeeccd46a2b8" {
		token, err := svc.GetJWTToken(creds.Username)
		if err != nil {
			svc.ResponseInternalError(w, err.Error())
			return
		}

		var response svc.ResGetJWT

		response.Status = true
		response.Code = http.StatusOK
		response.Message = "Success"
		response.Data.Token = token

		svc.ResponseWrite(w, response.Code, response)
	} else {
		svc.ResponseBadRequest(w, "Invalid authorization")
		return
	}
}
