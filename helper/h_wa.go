package helper

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"os"
	"strings"
	"time"

	whatsApp "github.com/dimaskiddo/whatsapp-go-mod"
	qrCode "github.com/skip2/go-qrcode"
)

func WAInit(timeout time.Duration) (*whatsApp.Conn, error) {
	conn, err := whatsApp.NewConn(timeout * time.Second)
	if err != nil {
		return nil, err
	}
	conn.SetClientName("WhatsApp Go", "WhatsApp Go")

	return conn, nil
}

func WASessionLoad(file string) (whatsApp.Session, error) {
	session := whatsApp.Session{}

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

func WASessionSave(file string, session whatsApp.Session) error {
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

func WASessionLogin(conn *whatsApp.Conn, file string, qr chan<- string) error {
	if conn != nil {
		_, err := os.Stat(file)
		if err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}

		session, err := conn.Login(qr)
		if err != nil {
			switch err.Error() {
			case "already logged in":
				return nil
			default:
				conn.EndConn()
				return errors.New("session not valid")
			}
		}

		err = WASessionSave(file, session)
		if err != nil {
			return err
		}
	} else {
		return errors.New("connection not valid")
	}

	return nil
}

func WASessionRestore(conn *whatsApp.Conn, sess whatsApp.Session, file string) error {
	if conn != nil {
		session, err := conn.RestoreSession(sess)
		if err != nil {
			switch err.Error() {
			case "already logged in":
				return nil
			default:
				err := conn.Logout()
				if err != nil {
					return err
				}

				conn.EndConn()
				return errors.New("session not valid")
			}
		}

		err = WASessionSave(file, session)
		if err != nil {
			return err
		}
	} else {
		return errors.New("connection not valid")
	}

	return nil
}

func WASessionLogout(conn *whatsApp.Conn, file string) error {
	if conn != nil {
		defer conn.EndConn()

		err := conn.Logout()
		if err != nil {
			return err
		}

		_, err = os.Stat(file)
		if err == nil {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("connection not valid")
	}

	return nil
}

func WAConnect(conn *whatsApp.Conn, timeout time.Duration, file string, qrcode chan<- string, errmsg chan<- error) {
	if conn != nil {
		defer conn.EndConn()

		chanqr := make(chan string)
		go func() {
			select {
			case tmp := <-chanqr:
				image, errImage := qrCode.New(tmp, qrCode.Medium)
				if errImage != nil {
					errmsg <- errImage
					return
				}

				png, errPNG := image.PNG(256)
				if errPNG != nil {
					errmsg <- errPNG
					return
				}

				qrcode <- base64.StdEncoding.EncodeToString(png)
			case <-time.After(timeout * time.Second):
				errmsg <- errors.New("qr code generate timeout")
			}
		}()

		session, err := WASessionLoad(file)
		if err != nil {
			err = WASessionLogin(conn, file, chanqr)
			if err != nil {
				errmsg <- err
				return
			}
		} else {
			err = WASessionRestore(conn, session, file)
			if err != nil {
				conn, err := WAInit(timeout)
				if err != nil {
					errmsg <- err
					return
				}

				err = WASessionLogin(conn, file, chanqr)
				if err != nil {
					errmsg <- err
					return
				}
			}
		}
	} else {
		errmsg <- errors.New("connection not valid")
		return
	}

	errmsg <- errors.New("")
	return
}

func WAMessageText(conn *whatsApp.Conn, jidDest string, msgText string, msgDelay time.Duration) error {
	if conn != nil {
		defer conn.EndConn()

		jidPrefix := "@s.whatsapp.net"
		if len(strings.SplitN(jidDest, "-", 2)) == 2 {
			jidPrefix = "@g.us"
		}

		content := whatsApp.TextMessage{
			Info: whatsApp.MessageInfo{
				RemoteJid: jidDest + jidPrefix,
			},
			Text: msgText,
		}

		<-time.After(msgDelay * time.Second)

		err := conn.Send(content)
		if err != nil {
			switch err.Error() {
			case "sending message timed out":
				return nil
			default:
				return err
			}
		}
	} else {
		return errors.New("connection not valid")
	}

	return nil
}
