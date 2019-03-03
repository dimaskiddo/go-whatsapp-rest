package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	hlp "github.com/dimaskiddo/go-whatsapp-rest/helper"
	svc "github.com/dimaskiddo/go-whatsapp-rest/service"
)

type reqWhatsAppLogin struct {
	Format  string `json:"format"`
	Timeout int    `json:"timeout"`
}

type resWhatsAppLogin struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		QRCode  string `json:"qrcode"`
		Timeout int    `json:"timeout"`
	} `json:"data"`
}

type reqWhatsAppSendMessageText struct {
	MSISDN  string `json:"msisdn"`
	Message string `json:"message"`
	Delay   int    `json:"delay"`
}

func WhatsAppLogin(w http.ResponseWriter, r *http.Request) {
	jid, err := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppLogin
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if reqBody.Timeout == 0 {
		reqBody.Timeout = 10
	}

	if len(reqBody.Format) == 0 {
		reqBody.Format = "json"
	}

	err = hlp.WAInit(jid, reqBody.Timeout)
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	file := svc.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	qrstr := make(chan string)
	errmsg := make(chan error)

	go func() {
		hlp.WAConnect(jid, reqBody.Timeout, file, qrstr, errmsg)
	}()

	select {
	case qrcode := <-qrstr:
		qrcode = "data:image/png;base64," + qrcode

		switch strings.ToLower(reqBody.Format) {
		case "json":
			var response resWhatsAppLogin

			response.Status = true
			response.Code = 200
			response.Message = "Success"
			response.Data.QRCode = qrcode
			response.Data.Timeout = reqBody.Timeout

			svc.ResponseWrite(w, response.Code, response)
		case "html":
			var response string

			response = `
        <html>
          <head>
            <title>WhatsApp Login</title>
          </head>
          <body>
              <img src="` + qrcode + `" />              
              <p>
                <b>Scan QR Code</b><br/>
                Timeout in ` + strconv.Itoa(reqBody.Timeout) + ` Second(s)
              </p>
            </center>
          </body>
        </html>
      `

			w.Write([]byte(response))
		default:
			svc.ResponseBadRequest(w, "")
		}
	case err := <-errmsg:
		if len(err.Error()) != 0 {
			svc.ResponseInternalError(w, err.Error())
			return
		}

		svc.ResponseSuccess(w, "")
	}
}

func WhatsAppLogout(w http.ResponseWriter, r *http.Request) {
	jid, err := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	file := svc.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	err = hlp.WASessionLogout(jid, file)
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	svc.ResponseSuccess(w, "")
}

func WhatsAppSendText(w http.ResponseWriter, r *http.Request) {
	jid, err := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessageText
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if len(reqBody.MSISDN) == 0 || len(reqBody.Message) == 0 {
		svc.ResponseBadRequest(w, "")
		return
	}

	err = hlp.WAMessageText(jid, reqBody.MSISDN, reqBody.Message, reqBody.Delay)
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	svc.ResponseSuccess(w, "")
}
