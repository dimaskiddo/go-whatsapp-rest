package index

import (
	"encoding/json"
	"net/http"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/router"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/server"
)

// GetAuth Function to Get Authorization Token
func GetAuth(w http.ResponseWriter, r *http.Request) {
	var reqBody auth.ReqGetBasic

	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if len(reqBody.Username) == 0 || len(reqBody.Password) == 0 {
		router.ResponseBadRequest(w, "invalid authorization")
		return
	}

	if reqBody.Password != server.Config.GetString("AUTH_BASIC_PASSWORD") {
		router.ResponseBadRequest(w, "invalid authorization")
		return
	}

	token, err := auth.GetJWTToken(reqBody.Username)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var response auth.ResGetJWT

	response.Status = true
	response.Code = http.StatusOK
	response.Message = "Success"
	response.Data.Token = token

	router.ResponseWrite(w, response.Code, response)
}
