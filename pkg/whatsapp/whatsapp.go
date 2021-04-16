package whatsapp

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"
	"sync"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	waproto "github.com/Rhymen/go-whatsapp/binary/proto"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/log"
)

var wac = make(map[string]*whatsapp.Conn)
var wacMutex = make(map[string]*sync.Mutex)

type waHandler struct {
	SessionConn   *whatsapp.Conn
	SessionJid    string
	SessionID     string
	SessionFile   string
	SessionStart  uint64
	ReconnectTime int
}

func (wah *waHandler) HandleError(err error) {
	_, eMatchFailed := err.(*whatsapp.ErrConnectionFailed)
	_, eMatchClosed := err.(*whatsapp.ErrConnectionClosed)

	if eMatchFailed || eMatchClosed {
		if WASessionExist(wah.SessionFile) && wah.SessionConn != nil {
			log.Println(log.LogLevelWarn, "whatsapp-error-handler", fmt.Sprintf("id %s connection closed unexpetedly, reconnecting after %d seconds", wah.SessionID, wah.ReconnectTime))

			<-time.After(time.Duration(wah.ReconnectTime) * time.Second)

			err := wah.SessionConn.Restore()
			if err != nil {
				log.Println(log.LogLevelError, "whatsapp-error-handler", fmt.Sprintf("id %s session restore error, %s", wah.SessionID, err.Error()))

				if WASessionExist(wah.SessionFile) {
					err := os.Remove(wah.SessionFile)
					if err != nil {
						log.Println(log.LogLevelError, "whatsapp-error-handler", fmt.Sprintf("id %s remove session file error, %s", wah.SessionID, err.Error()))
					}
				}

				_, _ = wah.SessionConn.Disconnect()
				wah.SessionConn.RemoveHandlers()

				delete(wac, wah.SessionJid)
			} else {
				wah.SessionStart = uint64(time.Now().Unix())
			}
		}
	} else {
		log.Println(log.LogLevelError, "whatsapp-error-handler", fmt.Sprintf("id %s caught an error, %s", wah.SessionID, err.Error()))
	}
}

func WAParseJID(jid string) string {
	components := strings.Split(jid, "@")

	if len(components) > 1 {
		jid = components[0]
	}

	suffix := "@s.whatsapp.net"

	if len(strings.SplitN(jid, "-", 2)) == 2 {
		suffix = "@g.us"
	}

	return jid + suffix
}

func WAGetSendMutexSleep() time.Duration {
	rand.Seed(time.Now().UnixNano())

	waitMin := 1000
	waitMax := 3000

	return time.Duration(rand.Intn(waitMax-rand.Intn(waitMin)) + waitMin)
}

func WASendWithMutex(jid string, content interface{}) (string, error) {
	mutex, ok := wacMutex[jid]

	if !ok {
		mutex = &sync.Mutex{}
		wacMutex[jid] = mutex
	}

	mutex.Lock()
	time.Sleep(WAGetSendMutexSleep() * time.Millisecond)

	id, err := wac[jid].Send(content)
	mutex.Unlock()

	return id, err
}

func WASyncVersion(conn *whatsapp.Conn, versionClientMajor int, versionClientMinor int, versionClientBuild int) (string, error) {
	// Bug Happend When Using This Function
	// Then Set Manualy WhatsApp Client Version
	// versionServer, err := whatsapp.CheckCurrentServerVersion()
	// if err != nil {
	// 	return "", err
	// }

	conn.SetClientVersion(versionClientMajor, versionClientMinor, versionClientBuild)
	versionClient := conn.GetClientVersion()

	return fmt.Sprintf("whatsapp version %v.%v.%v", versionClient[0], versionClient[1], versionClient[2]), nil
}

func WATestPing(conn *whatsapp.Conn) error {
	ok, err := conn.AdminTest()
	if !ok {
		if err != nil {
			return err
		}

		return errors.New("something when wrong while trying to ping, please check phone connectivity")
	}

	return nil
}

func WAGenerateQR(timeout int, chanqr chan string, qrstr chan<- string) {
	select {
	case tmp := <-chanqr:
		png, _ := qrcode.Encode(tmp, qrcode.Medium, 256)
		qrstr <- base64.StdEncoding.EncodeToString(png)
	}
}

func WASessionInit(jid string, versionClientMajor int, versionClientMinor int, versionClientBuild int, timeout int) error {
	if wac[jid] == nil {
		conn, err := whatsapp.NewConn(time.Duration(timeout) * time.Second)
		if err != nil {
			return err
		}
		conn.SetClientName("Go WhatsApp REST", "Go WhatsApp", "1.0")

		info, err := WASyncVersion(conn, versionClientMajor, versionClientMinor, versionClientBuild)
		if err != nil {
			return err
		}
		log.Println(log.LogLevelInfo, "whatsapp-session-init", info)

		wacMutex[jid] = &sync.Mutex{}
		wac[jid] = conn
	}

	return nil
}

func WASessionLoad(file string) (whatsapp.Session, error) {
	session := whatsapp.Session{}

	buffer, err := os.Open(file)
	if err != nil {
		return session, err
	}
	defer buffer.Close()

	err = gob.NewDecoder(buffer).Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func WASessionSave(file string, session whatsapp.Session) error {
	buffer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer buffer.Close()

	err = gob.NewEncoder(buffer).Encode(session)
	if err != nil {
		return err
	}

	return nil
}

func WASessionExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	return true
}

func WASessionConnect(jid string, versionClientMajor int, versionClientMinor int, versionClientBuild int, timeout int, file string, reconnect int, qrstr chan<- string, errmsg chan<- error) {
	chanqr := make(chan string)

	session, err := WASessionLoad(file)
	if err != nil {
		go func() {
			WAGenerateQR(timeout, chanqr, qrstr)
		}()

		err = WASessionLogin(jid, versionClientMajor, versionClientMinor, versionClientBuild, timeout, file, chanqr)
		if err != nil {
			errmsg <- err
			return
		}
		return
	}

	err = WASessionRestore(jid, versionClientMajor, versionClientMinor, versionClientBuild, timeout, file, session)
	if err != nil {
		go func() {
			WAGenerateQR(timeout, chanqr, qrstr)
		}()

		err = WASessionLogin(jid, versionClientMajor, versionClientMinor, versionClientBuild, timeout, file, chanqr)
		if err != nil {
			errmsg <- err
			return
		}
	}

	err = WATestPing(wac[jid])
	if err != nil {
		errmsg <- err
		return
	}

	wac[jid].AddHandler(&waHandler{
		SessionConn:   wac[jid],
		SessionJid:    jid,
		SessionID:     jid[0:len(jid)-4] + "xxxx",
		SessionFile:   file,
		SessionStart:  uint64(time.Now().Unix()),
		ReconnectTime: reconnect,
	})

	errmsg <- errors.New("")
	return
}

func WASessionLogin(jid string, versionClientMajor int, versionClientMinor int, versionClientBuild int, timeout int, file string, qrstr chan<- string) error {
	if wac[jid] != nil {
		if WASessionExist(file) {
			err := os.Remove(file)
			if err != nil {
				return err
			}
		}

		delete(wac, jid)
	}

	err := WASessionInit(jid, versionClientMajor, versionClientMinor, versionClientBuild, timeout)
	if err != nil {
		return err
	}

	session, err := wac[jid].Login(qrstr)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "already logged in":
			return nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return errors.New("connection is invalid")
		default:
			delete(wac, jid)
			return err
		}
	}

	err = WASessionSave(file, session)
	if err != nil {
		return err
	}

	return nil
}

func WASessionRestore(jid string, versionClientMajor int, versionClientMinor int, versionClientBuild int, timeout int, file string, sess whatsapp.Session) error {
	if wac[jid] != nil {
		if WASessionExist(file) {
			err := os.Remove(file)
			if err != nil {
				return err
			}
		}

		delete(wac, jid)
	}

	err := WASessionInit(jid, versionClientMajor, versionClientMinor, versionClientBuild, timeout)
	if err != nil {
		return err
	}

	session, err := wac[jid].RestoreWithSession(sess)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "already logged in":
			return nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return errors.New("connection is invalid")
		default:
			delete(wac, jid)
			return err
		}
	}

	err = WASessionSave(file, session)
	if err != nil {
		return err
	}

	return nil
}

func WASessionLogout(jid string, file string) error {
	if wac[jid] != nil {
		defer func() {
			_, _ = wac[jid].Disconnect()
			delete(wac, jid)
		}()

		wac[jid].RemoveHandlers()

		err := wac[jid].Logout()
		if err != nil {
			return err
		}

		if WASessionExist(file) {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("connection is invalid")
	}

	return nil
}

func WASessionValidate(jid string) error {
	if wac[jid] == nil {
		return errors.New("connection is invalid")
	}

	return nil
}

func WAMessageText(jid string, jidDest string, msgText string, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		Text: msgText,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}

func WAMessageDocument(jid string, jidDest string, msgDocumentStream multipart.File, msgDocumentType string, msgDocumentInfo string, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.DocumentMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		Content:  msgDocumentStream,
		Type:     msgDocumentType,
		FileName: msgDocumentInfo,
		Title:    msgDocumentInfo,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}

func WAMessageAudio(jid string, jidDest string, msgAudioStream multipart.File, msgAudioType string, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.AudioMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		Content: msgAudioStream,
		Type:    msgAudioType,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}

func WAMessageImage(jid string, jidDest string, msgImageStream multipart.File, msgImageType string, msgImageInfo string, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		Content: msgImageStream,
		Type:    msgImageType,
		Caption: msgImageInfo,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}

func WAMessageVideo(jid string, jidDest string, msgVideoStream multipart.File, msgVideoType string, msgVideoInfo string, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.VideoMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		Content: msgVideoStream,
		Type:    msgVideoType,
		Caption: msgVideoInfo,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}

func WAMessageLocation(jid string, jidDest string, msgLatitude float64, msgLongitude float64, msgQuotedID string, msgQuoted string) (string, error) {
	var id string

	err := WASessionValidate(jid)
	if err != nil {
		return "", errors.New(err.Error())
	}

	rJid := WAParseJID(jidDest)

	content := whatsapp.LocationMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: rJid,
		},
		DegreesLatitude:  msgLatitude,
		DegreesLongitude: msgLongitude,
	}

	if len(msgQuotedID) != 0 {
		msgQuotedProto := waproto.Message{
			Conversation: &msgQuoted,
		}

		ctxQuotedInfo := whatsapp.ContextInfo{
			QuotedMessageID: msgQuotedID,
			QuotedMessage:   &msgQuotedProto,
			Participant:     rJid,
		}

		content.ContextInfo = ctxQuotedInfo
	}

	id, err = WASendWithMutex(jid, content)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "sending message timed out":
			return id, nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return "", errors.New("connection is invalid")
		default:
			return "", err
		}
	}

	return id, nil
}
