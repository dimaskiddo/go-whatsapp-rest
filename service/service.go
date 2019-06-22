package service

// Initialize Function in Service
func init() {
	// Initialize Logger
	logInit()

	// Initialize Configuration
	configInit()

	// Initialize Cryptography
	cryptInit()

	// Initialize Router
	routerInit()
}
