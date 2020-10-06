package internal

import (
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/auth"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/router"

	"github.com/dimaskiddo/go-whatsapp-rest/internal/index"
	"github.com/dimaskiddo/go-whatsapp-rest/internal/whatsapp"
)

// LoadRoutes to Load Routes to Router
func LoadRoutes() {
	// Set Endpoint for Root Functions
	router.Router.Get(router.RouterBasePath, index.GetIndex)
	router.Router.Get(router.RouterBasePath+"/health", index.GetHealth)

	// Set Endpoint for Authorization Functions
	router.Router.With(auth.Basic).Get(router.RouterBasePath+"/auth", index.GetAuth)

	// Set Endpoint for WhatsApp Functions
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/login", whatsapp.WhatsAppLogin)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/text", whatsapp.WhatsAppSendText)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/location", whatsapp.WhatsAppSendLocation)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/document", whatsapp.WhatsAppSendDocument)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/audio", whatsapp.WhatsAppSendAudio)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/image", whatsapp.WhatsAppSendImage)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/send/video", whatsapp.WhatsAppSendVideo)
	router.Router.With(auth.JWT).Post(router.RouterBasePath+"/logout", whatsapp.WhatsAppLogout)
}
