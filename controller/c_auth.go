package controller

import (
	"encoding/json"
	"net/http"

	svc "github.com/dimaskiddo/whatsapp-go-rest/service"
)

// GetAuth Function to Get Authorization Token
func GetAuth(w http.ResponseWriter, r *http.Request) {
	var reqBody svc.ReqGetBasic
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if len(reqBody.Username) == 0 || len(reqBody.Password) == 0 {
		svc.ResponseBadRequest(w, "Invalid authorization")
		return
	}

	if reqBody.Password == "83e4060e-78e1-4fe5-9977-aeeccd46a2b8" {
		token, err := svc.GetJWTToken(reqBody.Username)
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
