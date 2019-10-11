package ctl

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/auth"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/libs"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/router"
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
	MSISDN        string `json:"msisdn"`
	Message       string `json:"message"`
	QuotedID      string `json:"quoteid"`
	QuotedMessage string `json:"quotedmsg"`
	Delay         int    `json:"delay"`
}

type resWhatsAppSendMessage struct {
	MessageID string `json:"msgid"`
}

func WhatsAppLogin(w http.ResponseWriter, r *http.Request) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppLogin
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if len(reqBody.Output) == 0 {
		reqBody.Output = "json"
	}

	if reqBody.Timeout == 0 {
		reqBody.Timeout = 5
	}

	file := hlp.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	qrstr := make(chan string)
	errmsg := make(chan error)

	go func() {
		libs.WASessionConnect(jid, reqBody.Timeout, file, qrstr, errmsg)
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

			router.ResponseWrite(w, response.Code, response)
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
			router.ResponseBadRequest(w, "")
		}
	case err := <-errmsg:
		if len(err.Error()) != 0 {
			router.ResponseInternalError(w, err.Error())
			return
		}

		router.ResponseSuccess(w, "")
	}
}

func WhatsAppLogout(w http.ResponseWriter, r *http.Request) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	file := hlp.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	err = libs.WASessionLogout(jid, file)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	router.ResponseSuccess(w, "")
}

func WhatsAppSendText(w http.ResponseWriter, r *http.Request) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessage
	_ = json.NewDecoder(r.Body).Decode(&reqBody)

	if len(reqBody.MSISDN) == 0 || len(reqBody.Message) == 0 {
		router.ResponseBadRequest(w, "")
		return
	}

	id, err := libs.WAMessageText(jid, reqBody.MSISDN, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage, reqBody.Delay)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var resBody resWhatsAppSendMessage
	resBody.MessageID = id

	router.ResponseSuccessWithData(w, "", resBody)
}

func WhatsAppSendImage(w http.ResponseWriter, r *http.Request) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	err = r.ParseMultipartForm(hlp.Config.GetInt64("SERVER_UPLOAD_LIMIT"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessage

	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.Message = r.FormValue("message")
	reqBody.QuotedID = r.FormValue("qoutedid")
	reqBody.QuotedMessage = r.FormValue("qoutedmsg")
	reqDelay := r.FormValue("delay")

	if len(reqDelay) == 0 {
		reqBody.Delay = 0
	} else {
		reqBody.Delay, err = strconv.Atoi(reqDelay)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	mpFileStream, mpFileHeader, err := r.FormFile("image")
	if err != nil {
		router.ResponseBadRequest(w, err.Error())
		return
	}
	defer mpFileStream.Close()

	mpFileType := mpFileHeader.Header.Get("Content-Type")

	if len(reqBody.MSISDN) == 0 || len(reqBody.Message) == 0 {
		router.ResponseBadRequest(w, "")
		return
	}

	id, err := libs.WAMessageImage(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage, reqBody.Delay)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var resBody resWhatsAppSendMessage
	resBody.MessageID = id

	router.ResponseSuccessWithData(w, "", resBody)
}
