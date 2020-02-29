package main

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/libs"
)

// Initialize function to restore sessions stored
func init() {
	files, err := filepath.Glob(hlp.Config.GetString("SERVER_STORE_PATH") + "/*.gob")
	if err == nil {
		for _, file := range files {
			jid := strings.TrimSuffix(filepath.Base(file), path.Ext(file))
			qrstr := make(chan string)
			errmsg := make(chan error)

			hlp.LogPrintln(hlp.LogLevelInfo, "whatsapp", "Attempting to restore session state for: "+jid)

			go func() {
				libs.WASessionConnect(jid, 20, file, qrstr, errmsg)
			}()

			select {
			case err := <-errmsg:
				if len(err.Error()) != 0 {
					hlp.LogPrintln(hlp.LogLevelWarn, "whatsapp", "Failed to restore session: "+err.Error())
					return
				}

				hlp.LogPrintln(hlp.LogLevelInfo, "whatsapp", "Session restored: "+jid)
			}
		}
	}
}
