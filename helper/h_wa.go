package helper

import (
	"encoding/gob"
	"errors"
	"os"
	"strings"
	"time"

	whatsApp "github.com/Rhymen/go-whatsapp"
	qrCode "github.com/skip2/go-qrcode"
)

// WhatsAppConnection Map Variable
var WhatsAppConnection = make(map[string]*whatsApp.Conn)

// WhatsAppConnect Function
func WhatsAppConnect(msisdn string, timeout time.Duration) error {
	// Remove Exsisiting Connection If Exist
	if WhatsAppConnection[msisdn] == nil {
		var errConnectionCreate error

		// Create New Connection
		WhatsAppConnection[msisdn], errConnectionCreate = whatsApp.NewConn(timeout * time.Second)
		if errConnectionCreate != nil {
			// Return Connection Creation Error
			return errConnectionCreate
		}

		// Set Connection Identifier
		WhatsAppConnection[msisdn].SetClientName("WhatsApp Go", "WhatsApp Go")
	}

	// Return No Error
	return nil
}

// WhatsAppLogin Function
func WhatsAppLogin(msisdn string, timeout time.Duration, fileSession string, chanQRCode chan<- []byte, chanError chan<- error) {
	// Check If Connection Exist
	if WhatsAppConnection[msisdn] != nil {
		// Try To Restore Session From Session File
		errSessionRestore := WhatsAppSessionRestore(msisdn, fileSession)
		if errSessionRestore != nil {
			// If Restore Session Error Then
			// Check If Session File Exist
			_, errFileSessionExist := os.Stat(fileSession)
			if errFileSessionExist == nil {
				// If Session File Exist
				// Try To Remove Session File
				errFileSessionRemove := os.Remove(fileSession)
				if errFileSessionRemove != nil {
					// Return Session File Remove Error
					// Using Error Message Channel
					chanError <- errFileSessionRemove
				}
			}

			// Create QR Code Data Channel With Type String
			qrCodeData := make(chan string)

			// Go Routine To Generata QR Code Data In Bytes
			go func() {
				select {
				case qrCodeString := <-qrCodeData:
					// If QR Code Data Channel Got Data From Login Function Then
					// Create QR Code Image Data
					qrCodeImage, errQRCodeImage := qrCode.New(qrCodeString, qrCode.Medium)
					if errQRCodeImage != nil {
						// Return QR Code Image Creation Error
						// Using Error Message Channel
						chanError <- errQRCodeImage
						return
					}

					// Create QR Code Image PNG
					qrCodeImagePNG, errQRCodeImagePNG := qrCodeImage.PNG(256)
					if errQRCodeImagePNG != nil {
						// Return QR Code Image Creation Error
						// Using Error Message Channel
						chanError <- errQRCodeImagePNG
						return
					}

					// Return QR Code Image PNG Using QR Code Channel
					chanQRCode <- qrCodeImagePNG
				case <-time.After(timeout * time.Second):
					// If Got Timeout It Can Be Mean No QR Code Data Return By Login Function
					// Return QR Code Creation Timeout Error Using Error Message Channel
					chanError <- errors.New("qr code creation timeout")
				}
			}()

			// Try To Logging In Using Exsisting Connection
			dataSession, errLogin := WhatsAppConnection[msisdn].Login(qrCodeData)
			if errLogin != nil {
				// If Login Failed Return Some Error Value
				switch errLogin.Error() {
				case "already logged in":
					// If Login Failed Caused
					// By Already Logged In Then Return Empty String Error
					// Using Error Message Channel
					chanError <- errors.New("")
					return
				default:
					// Return Login Error
					// Using Error Message Channel
					chanError <- errLogin
					return
				}
			}

			// Try To Save Session Data To Session File
			errSessionSave := WhatsAppSessionSave(fileSession, dataSession)
			if errSessionSave != nil {
				// Return Session File Save Error
				// Using Error Message Channel
				chanError <- errSessionSave
				return
			}
		}
	} else {
		chanError <- errors.New("connection not found")
		return
	}

	// Return Empty String Error
	// Using Error Message Channel
	chanError <- errors.New("")
}

// WhatsAppLogout Function
func WhatsAppLogout(msisdn string, fileSession string) error {
	// Check If Connection Exist
	if WhatsAppConnection[msisdn] != nil {
		// Logout Connection And Make Session Invalidated
		errLogout := WhatsAppConnection[msisdn].Logout()
		if errLogout != nil {
			// Return Connection Logout Error
			return errLogout
		}

		// Check If Session File Exist
		_, errFileSessionExist := os.Stat(fileSession)
		if errFileSessionExist == nil {
			// If Session File Exist
			// Try To Remove Session File
			errFileSessionRemove := os.Remove(fileSession)
			if errFileSessionRemove != nil {
				// Return Session File Remove Error
				return errFileSessionRemove
			}
		}

		// Remove Connection From Map
		delete(WhatsAppConnection, msisdn)
	} else {
		// Return Connection Not Found Error
		return errors.New("connection not found")
	}

	// Return No Error
	return nil
}

// WhatsAppSessionRestore Function
func WhatsAppSessionRestore(msisdn string, fileSession string) error {
	// Check If Connection Exist
	if WhatsAppConnection[msisdn] != nil {
		// Create Empty Session Data As Session Comparator
		nilSession := whatsApp.Session{}

		// Try To Load Session Data From Session File
		dataSession, errSessionLoad := WhatsAppSessionLoad(fileSession)
		if errSessionLoad != nil {
			// Return Session File Load Error
			return errSessionLoad
		}

		// Check If Loaded Session Data Not The Same With Empty Session Data
		if dataSession.ClientId != nilSession.ClientId {
			// If Session Data Valid Then
			// Try To Restore Session Data In To Connection
			dataSession, errSessionRestore := WhatsAppConnection[msisdn].RestoreSession(dataSession)
			if errSessionRestore != nil {
				// If Restore Session Failed Return Some Error Value
				switch errSessionRestore.Error() {
				case "already logged in":
					// If Restore Session Failed Caused
					// By Already Logged In Then Return No Error
					return nil
				default:
					// Return Restore Session Error
					return errSessionRestore
				}
			}

			// Try To Re-Save Session Data From Restored Session To Session File
			errSessionSave := WhatsAppSessionSave(fileSession, dataSession)
			if errSessionSave != nil {
				// Return Session File Save Error
				return errSessionSave
			}
		} else {
			// Return Session Data Not Valid
			return errors.New("session data not valid")
		}
	} else {
		// Return Connection Not Found Error
		return errors.New("connection not found")
	}

	// Return No Error
	return nil
}

// WhatsAppSessionLoad Function
func WhatsAppSessionLoad(fileSession string) (whatsApp.Session, error) {
	// Create Empty Session Data
	dataSession := whatsApp.Session{}

	// Try To Load Session File
	fileSessionLoad, errFileSessionLoad := os.Open(fileSession)
	if errFileSessionLoad != nil {
		// Return Empty Session Data And Session File Open Error
		return dataSession, errFileSessionLoad
	}

	// Close Session File When Function Is Done
	defer fileSessionLoad.Close()

	// Try To Decode Session File Content And Restore It Session Data
	errSesionDecode := gob.NewDecoder(fileSessionLoad).Decode(&dataSession)
	if errSesionDecode != nil {
		// Return Empty Session Data And Session File Content Decode Error
		return dataSession, errSesionDecode
	}

	// Return Session Data And No Error
	return dataSession, nil
}

// WhatsAppSessionSave Function
func WhatsAppSessionSave(fileSession string, dataSession whatsApp.Session) error {
	// Try To Create Session File
	fileSessionSave, errFileSessionSave := os.Create(fileSession)
	if errFileSessionSave != nil {
		// Return Session File Create Error
		return errFileSessionSave
	}

	// Close Session File When Function Is Done
	defer fileSessionSave.Close()

	// Try To Encode Session Data Content And Save It Session File
	errSessionEncode := gob.NewEncoder(fileSessionSave).Encode(dataSession)
	if errSessionEncode != nil {
		// Return Session Data Encode Error
		return errSessionEncode
	}

	// Return No Error
	return nil
}

// WhatsAppSendMessageText Function
func WhatsAppSendMessageText(msisdn string, fileSession string, msisdnDestination string, messageText string, messageDelay time.Duration) error {
	// Check If Connection Exist
	if WhatsAppConnection[msisdn] != nil {
		// Try To Restore Session From Session File
		errSessionRestore := WhatsAppSessionRestore(msisdn, fileSession)
		if errSessionRestore != nil {
			// Return Session Restore Error
			return errSessionRestore
		}

		// Set RemoteJID Prefix
		jidPrefix := "@s.whatsapp.net"
		jidDestinationCheck := strings.SplitN(msisdnDestination, "-", 2)
		if len(jidDestinationCheck) == 2 {
			jidPrefix = "@g.us"
		}

		// Set Message Text Content
		msgContent := whatsApp.TextMessage{
			Info: whatsApp.MessageInfo{
				RemoteJid: msisdnDestination + jidPrefix,
			},
			Text: messageText,
		}

		// Delay Before Send Message Text
		<-time.After(messageDelay * time.Second)

		// Try To Send Message Text
		errSendMessageText := WhatsAppConnection[msisdn].Send(msgContent)
		if errSendMessageText != nil {
			// Return Send Message Text Error
			return errSendMessageText
		}
	} else {
		// Return Connection Not Found Error
		return errors.New("connection not found")
	}

	// Return No Error
	return nil
}
