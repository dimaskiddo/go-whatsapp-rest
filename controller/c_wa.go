package controller

import (
	"net/http"
	"strconv"
	"strings"

	hlp "github.com/dimaskiddo/go-whatsapp-rest/helper"
	svc "github.com/dimaskiddo/go-whatsapp-rest/service"
)

type reqWhatsAppLogin struct {
	Output  string `json:"output"`
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

type reqWhatsAppSendMessage struct {
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

	err = r.ParseForm()
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppLogin

	reqBody.Output = r.FormValue("output")
	reqTimeout := r.FormValue("timeout")

	if len(reqBody.Output) == 0 {
		reqBody.Output = "json"
	}

	if len(reqTimeout) == 0 {
		reqBody.Timeout = 10
	} else {
		reqBody.Timeout, err = strconv.Atoi(reqTimeout)
		if err != nil {
			svc.ResponseInternalError(w, err.Error())
			return
		}
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

		switch strings.ToLower(reqBody.Output) {
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
              <b>QR Code Scan</b>
              <br/>
              Timeout in ` + strconv.Itoa(reqBody.Timeout) + ` Second(s)
            </p>
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

	err = r.ParseForm()
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessage

	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.Message = r.FormValue("message")
	reqDelay := r.FormValue("delay")

	if len(reqDelay) == 0 {
		reqBody.Delay = 0
	} else {
		reqBody.Delay, err = strconv.Atoi(reqDelay)
		if err != nil {
			svc.ResponseInternalError(w, err.Error())
			return
		}
	}

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

func WhatsAppSendImage(w http.ResponseWriter, r *http.Request) {
	jid, err := svc.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	err = r.ParseMultipartForm(svc.Config.GetInt64("SERVER_UPLOAD_LIMIT"))
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessage

	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.Message = r.FormValue("message")
	reqDelay := r.FormValue("delay")

	if len(reqDelay) == 0 {
		reqBody.Delay = 0
	} else {
		reqBody.Delay, err = strconv.Atoi(reqDelay)
		if err != nil {
			svc.ResponseInternalError(w, err.Error())
			return
		}
	}

	mpFileStream, mpFileHeader, err := r.FormFile("image")
	if err != nil {
		svc.ResponseBadRequest(w, err.Error())
		return
	}
	defer mpFileStream.Close()

	mpFileType := mpFileHeader.Header.Get("Content-Type")

	if len(reqBody.MSISDN) == 0 || len(reqBody.Message) == 0 {
		svc.ResponseBadRequest(w, "")
		return
	}

	err = hlp.WAMessageImage(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.Message, reqBody.Delay)
	if err != nil {
		svc.ResponseInternalError(w, err.Error())
		return
	}

	svc.ResponseSuccess(w, "")
}
