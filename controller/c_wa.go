package controller

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	hlp "github.com/dimaskiddo/whatsapp-go-rest/helper"
	svc "github.com/dimaskiddo/whatsapp-go-rest/service"
)

type formatWhatsAppLogin struct {
	Format  string        `json:"format"`
	Timeout time.Duration `json:"timeout"`
}

type formatWhatsAppLoginQR struct {
	Status  bool              `json:"status"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

type formatWhatsAppSendMessageText struct {
	MSISDN  string        `json:"msisdn"`
	Message string        `json:"message"`
	Delay   time.Duration `json:"delay"`
}

func WhatsAppLogin(w http.ResponseWriter, r *http.Request) {
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		msisdn := string(dataClaims)

		var dataRequest formatWhatsAppLogin

		errDataRequest := json.NewDecoder(r.Body).Decode(&dataRequest)
		if errDataRequest != nil {
			svc.ResponseInternalError(w, errDataRequest.Error())
		} else {
			if len(dataRequest.Format) != 0 {
				errConnectionCreate := hlp.WhatsAppConnect(msisdn, dataRequest.Timeout)
				if errConnectionCreate != nil {
					svc.ResponseInternalError(w, errConnectionCreate.Error())
				} else {
					fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

					loginQRCode := make(chan []byte)
					loginError := make(chan error)

					go func() {
						hlp.WhatsAppLogin(msisdn, dataRequest.Timeout, fileSession, loginQRCode, loginError)
					}()

					select {
					case qrCodeLogin := <-loginQRCode:
						qrCodeEncoded := base64.StdEncoding.EncodeToString(qrCodeLogin)

						switch strings.ToLower(dataRequest.Format) {
						case "json":
							var response formatWhatsAppLoginQR

							response.Status = true
							response.Code = 200
							response.Message = "Success"
							response.Data = map[string]string{
								"qrcode": "data:image/png;base64," + qrCodeEncoded,
							}

							svc.ResponseWrite(w, response.Code, response)
						case "html":
							var response string

							response = `
              <html>
                <head>
                  <title>WhatsApp Login</title>
                </head>

                <body>
                  <img src="data:image/png;base64,` + qrCodeEncoded + `" />
                </body>
              </html>
              `

							w.Write([]byte(response))
						default:
							svc.ResponseBadRequest(w, "")
						}
					case errLogin := <-loginError:
						if len(errLogin.Error()) != 0 {
							svc.ResponseInternalError(w, errLogin.Error())
						} else {
							svc.ResponseSuccess(w, "")
						}
					}
				}
			} else {
				svc.ResponseBadRequest(w, "")
			}
		}
	}
}

func WhatsAppLogout(w http.ResponseWriter, r *http.Request) {
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		msisdn := string(dataClaims)
		fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

		errLogout := hlp.WhatsAppLogout(msisdn, fileSession)
		if errLogout != nil {
			svc.ResponseInternalError(w, errLogout.Error())
		} else {
			svc.ResponseSuccess(w, "")
		}
	}
}

func WhatsAppSendMessageText(w http.ResponseWriter, r *http.Request) {
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		msisdn := string(dataClaims)

		var dataRequest formatWhatsAppSendMessageText

		errDataRequest := json.NewDecoder(r.Body).Decode(&dataRequest)
		if errDataRequest != nil {
			svc.ResponseInternalError(w, errDataRequest.Error())
		} else {
			fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

			errSendMessageText := hlp.WhatsAppSendMessageText(msisdn, fileSession, dataRequest.MSISDN, dataRequest.Message, dataRequest.Delay)
			if errSendMessageText != nil {
				svc.ResponseInternalError(w, errSendMessageText.Error())
			} else {
				svc.ResponseSuccess(w, "")
			}
		}
	}
}
