package main

import (
	ctl "github.com/dimaskiddo/whatsapp-go-rest/controller"
	svc "github.com/dimaskiddo/whatsapp-go-rest/service"
)

// Routes Initialization Function
func initRoutes() {
	// Set Endpoint for Root Functions
	svc.Router.HandleFunc(svc.RouterBasePath, ctl.GetIndex).Methods("GET")
	svc.Router.HandleFunc(svc.RouterBasePath+"health", ctl.GetHealth).Methods("GET")

	// Set Endpoint for Authorization Functions
	svc.Router.Handle(svc.RouterBasePath+"auth", svc.AuthBasic(ctl.GetAuth)).Methods("GET", "POST")

	// Set Endpoint for WhatsApp Functions
	svc.Router.Handle(svc.RouterBasePath+"login", svc.AuthJWT(ctl.WhatsAppLogin)).Methods("POST")
	svc.Router.Handle(svc.RouterBasePath+"messagetext", svc.AuthJWT(ctl.WhatsAppSendText)).Methods("POST")
	svc.Router.Handle(svc.RouterBasePath+"logout", svc.AuthJWT(ctl.WhatsAppLogout)).Methods("POST")
}
