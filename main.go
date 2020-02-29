package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/libs"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/router"
)

// Server Variable
var svr *hlp.Server

// Init Function
func init() {
	// Initialize Server
	svr = hlp.NewServer(router.Router)
	restore()
}

// Main Function
func main() {
	// Starting Server
	svr.Start()

	// Make Channel for OS Signal
	sig := make(chan os.Signal, 1)

	// Notify Any Signal to OS Signal Channel
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	// Return OS Signal Channel
	// As Exit Sign
	<-sig

	// Log Break Line
	fmt.Println("")

	// Stopping Server
	defer svr.Stop()
}

func restore() {
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
