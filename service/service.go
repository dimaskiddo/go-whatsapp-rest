package service

// Initialize Function in Utils
func Initialize() {
	// Initialize Logger
	logInit()

	// Initialize Configuration
	configInit()

	// Initialize Cryptography
	cryptInit()

	// Initialize Router
	routerInit()
}
