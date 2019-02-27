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

// formatWhatsAppLogin Struct
type formatWhatsAppLogin struct {
	Format  string        `json:"format"`
	Timeout time.Duration `json:"timeout"`
}

// formatWhatsAppLoginQR Struct
type formatWhatsAppLoginQR struct {
	Status  bool              `json:"status"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

// formatWhatsAppSendMessageText Struct
type formatWhatsAppSendMessageText struct {
	MSISDN  string        `json:"msisdn"`
	Message string        `json:"message"`
	Delay   time.Duration `json:"delay"`
}

// WhatsAppLogin Function
func WhatsAppLogin(w http.ResponseWriter, r *http.Request) {
	// Get Claims Data From JWT Authorization
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		// Response With Internal Error Claims Data
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		// Get MSISDN From Claims Data
		msisdn := string(dataClaims)

		// Create Request Data Variable
		var dataRequest formatWhatsAppLogin

		// Parse Request Body Data To Request Data Variable
		_ = json.NewDecoder(r.Body).Decode(&dataRequest)

		// Check If Request Data Has Format Option
		if len(dataRequest.Format) == 0 {
			// If Format Options Empty Then
			// Set To Default Option (JSON)
			dataRequest.Format = "json"
		}

		// Check If Request Data Has Timeout Option
		if dataRequest.Timeout == 0 {
			// If Timeout Options Empty Then
			// Set To Default Option (10 Seconds)
			dataRequest.Timeout = 10
		}

		// Try To Create Connection With MSISDN As Identifier
		errConnectionCreate := hlp.WhatsAppConnect(msisdn, dataRequest.Timeout)
		if errConnectionCreate != nil {
			// Response With Internal Error Create Connection
			svc.ResponseInternalError(w, errConnectionCreate.Error())
		} else {
			// Set Session File Variable To Server Store Path With MSISDN As Identifier
			fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

			// Create Channel For QR Code Login String
			loginQRCode := make(chan []byte)

			// Create Channel For Login Error
			loginError := make(chan error)

			// Routine For Logging In
			go func() {
				// Try To Login With Created Connrection
				hlp.WhatsAppLogin(msisdn, dataRequest.Timeout, fileSession, loginQRCode, loginError)
			}()

			select {
			case qrCodeLogin := <-loginQRCode:
				// If QR Code Login Channel Got Data From Login Function Then
				// Encode QR Code To Base64 Format
				qrCodeEncoded := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCodeLogin)

				// Response Encoded QR Code With
				// Selected Output Format Option
				switch strings.ToLower(dataRequest.Format) {
				case "json":
					// Create Response Format
					var response formatWhatsAppLoginQR

					// Fill in Response Data
					response.Status = true
					response.Code = 200
					response.Message = "Success"
					response.Data = map[string]string{
						"qrcode": qrCodeEncoded,
					}

					// Write Response Data
					svc.ResponseWrite(w, response.Code, response)
				case "html":
					// Create Response Content
					var response string

					// Fill in Response Content
					response = `
          <html>
            <head>
              <title>WhatsApp Login</title>
            </head>

            <body>
              <img src="` + qrCodeEncoded + `" />
            </body>
          </html>
          `

					// Write Response
					w.Write([]byte(response))
				default:
					// If Format Option Doesn't Match
					// Response With Bad Request
					svc.ResponseBadRequest(w, "")
				}
			case errLogin := <-loginError:
				// If Error Login Channel Got Data From Login Function Then
				// Check If Error Message Empty
				if len(errLogin.Error()) != 0 {
					// Response With Internal Error Login
					svc.ResponseInternalError(w, errLogin.Error())
				} else {
					// Response With Success
					svc.ResponseSuccess(w, "")
				}
			}
		}
	}
}

// WhatsAppLogout Function
func WhatsAppLogout(w http.ResponseWriter, r *http.Request) {
	// Get Claims Data From JWT Authorization
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		// Response With Internal Error Claims Data
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		// Get MSISDN From Claims Data
		msisdn := string(dataClaims)

		// Set Session File Variable To Server Store Path With MSISDN As Identifier
		fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

		// Try To Logout With Current Connection
		errLogout := hlp.WhatsAppLogout(msisdn, fileSession)
		if errLogout != nil {
			// Response With Internal Error Logout
			svc.ResponseInternalError(w, errLogout.Error())
		} else {
			// Response With Success
			svc.ResponseSuccess(w, "")
		}
	}
}

// WhatsAppSendMessageText Function
func WhatsAppSendMessageText(w http.ResponseWriter, r *http.Request) {
	// Get Claims Data From JWT Authorization
	dataClaims, errDataClaims := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if errDataClaims != nil {
		// Response With Internal Error Claims Data
		svc.ResponseInternalError(w, errDataClaims.Error())
	} else {
		// Get MSISDN From Claims Data
		msisdn := string(dataClaims)

		// Create Request Data Variable
		var dataRequest formatWhatsAppSendMessageText

		// Parse Request Body Data To Request Data Variable
		_ = json.NewDecoder(r.Body).Decode(&dataRequest)

		// Check If Request Data Has Destination MSISDN And Message Text
		if len(dataRequest.MSISDN) != 0 || len(dataRequest.Message) != 0 {
			// Set Session File Variable To Server Store Path With MSISDN As Identifier
			fileSession := svc.Config.GetString("SERVER_STORE_PATH") + "/" + msisdn + ".gob"

			// Try To Send Message Text
			errSendMessageText := hlp.WhatsAppSendMessageText(msisdn, fileSession, dataRequest.MSISDN, dataRequest.Message, dataRequest.Delay)
			if errSendMessageText != nil {
				// Response With Send Message Text Error
				svc.ResponseInternalError(w, errSendMessageText.Error())
			} else {
				// Response With Success
				svc.ResponseSuccess(w, "")
			}
		} else {
			// Response With Bad Request
			svc.ResponseBadRequest(w, "")
		}
	}
}
