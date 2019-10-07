package main

import (
	"github.com/dimaskiddo/go-whatsapp-rest/ctl"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/auth"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/router"
)

// Initialize Function in Main Route
func init() {
	// Set Endpoint for Root Functions
	router.Router.Get(router.RouterBasePath, ctl.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", ctl.GetHealth)

	// Set Endpoint for Authorization Functions
	router.Router.With(auth.Basic).Get(router.RouterBasePath+"/auth", ctl.GetAuth)

	// Set Endpoint for WhatsApp Functions
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/login", ctl.WhatsAppLogin)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/sendtext", ctl.WhatsAppSendText)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/sendimage", ctl.WhatsAppSendImage)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/logout", ctl.WhatsAppLogout)
}
