package service

// Initialize Function in Utils
func Initialize() {
	// Initialize Logger
	initLog()

	// Initialize Configuration
	initConfig()

	// Initialize Cryptography
	initCrypt()

	// Initialize Router
	initRouter()
}
