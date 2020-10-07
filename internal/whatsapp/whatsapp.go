package whatsapp

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/router"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/server"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/whatsapp"
)

type reqWhatsAppLogin struct {
	Output    string
	Reconnect int
	Timeout   int
	WhatsApp  struct {
		Client struct {
			Version struct {
				Major int
				Minor int
				Build int
			}
		}
	}
}

type resWhatsAppLogin struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		QRCode    string `json:"qrcode"`
		Reconnect int    `json:"reconnect"`
		Timeout   int    `json:"timeout"`
	} `json:"data"`
}

type reqWhatsAppSendMessage struct {
	MSISDN        string
	Message       string
	QuotedID      string
	QuotedMessage string
}

type reqWhatsAppSendLocation struct {
	MSISDN        string
	Latitude      float64
	Longitude     float64
	QuotedID      string
	QuotedMessage string
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

	r.ParseForm()

	var reqBody reqWhatsAppLogin
	reqBody.Output = r.FormValue("output")
	reqReconnect := r.FormValue("reconnect")
	reqTimeout := r.FormValue("timeout")

	reqVersionClientMajor := r.FormValue("client_version_major")
	reqVersionClientMinor := r.FormValue("client_version_minor")
	reqVersionClientBuild := r.FormValue("client_version_build")

	if len(reqBody.Output) == 0 {
		reqBody.Output = "json"
	}

	if len(reqReconnect) == 0 {
		reqBody.Reconnect = 30
	} else {
		reqBody.Reconnect, err = strconv.Atoi(reqReconnect)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	if len(reqTimeout) == 0 {
		reqBody.Timeout = 5
	} else {
		reqBody.Timeout, err = strconv.Atoi(reqTimeout)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	if len(reqVersionClientMajor) == 0 {
		reqBody.WhatsApp.Client.Version.Major = server.Config.GetInt("WHATSAPP_CLIENT_VERSION_MAJOR")
	} else {
		reqBody.WhatsApp.Client.Version.Major, err = strconv.Atoi(reqVersionClientMajor)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	if len(reqVersionClientMinor) == 0 {
		reqBody.WhatsApp.Client.Version.Minor = server.Config.GetInt("WHATSAPP_CLIENT_VERSION_MINOR")
	} else {
		reqBody.WhatsApp.Client.Version.Minor, err = strconv.Atoi(reqVersionClientMinor)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	if len(reqVersionClientBuild) == 0 {
		reqBody.WhatsApp.Client.Version.Build = server.Config.GetInt("WHATSAPP_CLIENT_VERSION_BUILD")
	} else {
		reqBody.WhatsApp.Client.Version.Build, err = strconv.Atoi(reqVersionClientBuild)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	file := server.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	qrstr := make(chan string)
	errmsg := make(chan error)

	go func() {
		whatsapp.WASessionConnect(jid, reqBody.WhatsApp.Client.Version.Major, reqBody.WhatsApp.Client.Version.Minor, reqBody.WhatsApp.Client.Version.Build, reqBody.Timeout, file, reqBody.Reconnect, qrstr, errmsg)
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

	file := server.Config.GetString("SERVER_STORE_PATH") + "/" + jid + ".gob"

	err = whatsapp.WASessionLogout(jid, file)
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

	r.ParseForm()

	var reqBody reqWhatsAppSendMessage
	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.Message = r.FormValue("message")
	reqBody.QuotedID = r.FormValue("quotedid")
	reqBody.QuotedMessage = r.FormValue("quotedmsg")

	if len(reqBody.MSISDN) == 0 || len(reqBody.Message) == 0 {
		router.ResponseBadRequest(w, "")
		return
	}

	id, err := whatsapp.WAMessageText(jid, reqBody.MSISDN, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var resBody resWhatsAppSendMessage
	resBody.MessageID = id

	router.ResponseSuccessWithData(w, "", resBody)
}

func WhatsAppSendContent(w http.ResponseWriter, r *http.Request, c string) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	err = r.ParseMultipartForm(server.Config.GetInt64("SERVER_UPLOAD_LIMIT"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var reqBody reqWhatsAppSendMessage
	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.QuotedID = r.FormValue("quotedid")
	reqBody.QuotedMessage = r.FormValue("quotedmsg")

	var mpFileStream multipart.File
	var mpFileHeader *multipart.FileHeader

	switch c {
	case "document":
		mpFileStream, mpFileHeader, err = r.FormFile("document")
		reqBody.Message = mpFileHeader.Filename

	case "audio":
		mpFileStream, mpFileHeader, err = r.FormFile("audio")

	case "image":
		mpFileStream, mpFileHeader, err = r.FormFile("image")
		reqBody.Message = r.FormValue("message")

	case "video":
		mpFileStream, mpFileHeader, err = r.FormFile("video")
		reqBody.Message = r.FormValue("message")
	}

	if err != nil {
		router.ResponseBadRequest(w, err.Error())
		return
	}
	defer mpFileStream.Close()

	mpFileType := mpFileHeader.Header.Get("Content-Type")

	if len(reqBody.MSISDN) == 0 {
		router.ResponseBadRequest(w, "")
		return
	}

	var id string

	switch c {
	case "document":
		id, err = whatsapp.WAMessageDocument(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

	case "audio":
		id, err = whatsapp.WAMessageAudio(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.QuotedID, reqBody.QuotedMessage)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

	case "image":
		id, err = whatsapp.WAMessageImage(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

	case "video":
		id, err = whatsapp.WAMessageVideo(jid, reqBody.MSISDN, mpFileStream, mpFileType, reqBody.Message, reqBody.QuotedID, reqBody.QuotedMessage)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}
	}

	var resBody resWhatsAppSendMessage
	resBody.MessageID = id

	router.ResponseSuccessWithData(w, "", resBody)
}

func WhatsAppSendDocument(w http.ResponseWriter, r *http.Request) {
	WhatsAppSendContent(w, r, "document")
}

func WhatsAppSendAudio(w http.ResponseWriter, r *http.Request) {
	WhatsAppSendContent(w, r, "audio")
}

func WhatsAppSendImage(w http.ResponseWriter, r *http.Request) {
	WhatsAppSendContent(w, r, "image")
}

func WhatsAppSendVideo(w http.ResponseWriter, r *http.Request) {
	WhatsAppSendContent(w, r, "video")
}

func WhatsAppSendLocation(w http.ResponseWriter, r *http.Request) {
	jid, err := auth.GetJWTClaims(r.Header.Get("X-JWT-Claims"))
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	r.ParseForm()

	var reqBody reqWhatsAppSendLocation
	reqBody.MSISDN = r.FormValue("msisdn")
	reqBody.QuotedID = r.FormValue("quotedid")
	reqBody.QuotedMessage = r.FormValue("quotedmsg")

	reqBody.Latitude, err = strconv.ParseFloat(r.FormValue("latitude"), 64)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	reqBody.Longitude, err = strconv.ParseFloat(r.FormValue("longitude"), 64)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	if len(reqBody.MSISDN) == 0 {
		router.ResponseBadRequest(w, "")
		return
	}

	id, err := whatsapp.WAMessageLocation(jid, reqBody.MSISDN, reqBody.Latitude, reqBody.Longitude, reqBody.QuotedID, reqBody.QuotedMessage)
	if err != nil {
		router.ResponseInternalError(w, err.Error())
		return
	}

	var resBody resWhatsAppSendMessage
	resBody.MessageID = id

	router.ResponseSuccessWithData(w, "", resBody)
}
