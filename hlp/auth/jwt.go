package auth

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
	"github.com/dimaskiddo/go-whatsapp-rest/hlp/router"
)

// ResGetJWT Struct
type ResGetJWT struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

// JWT Claims Data Struct
type jwtClaimsData struct {
	Data string `json:"data"`
	jwt.StandardClaims
}

// JWT Function as Midleware for JWT Authorization
func JWT(next http.Handler) http.Handler {
	// Return Next HTTP Handler Function, If Authorization is Valid
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse HTTP Header Authorization
		authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		// Check HTTP Header Authorization Section
		// Authorization Section Length Should Be 2
		// The First Authorization Section Should Be "Bearer"
		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			hlp.LogPrintln(hlp.LogLevelWarn, "http-access", "unauthorized method "+r.Method+" at URI "+r.RequestURI)
			router.ResponseUnauthorized(w)
			return
		}

		// The Second Authorization Section Should Be The Credentials Payload
		authPayload := authHeader[1]
		if len(authPayload) == 0 {
			router.ResponseBadRequest(w, "")
			return
		}

		// Get Authorization Claims From JWT Token
		authClaims, err := jwtClaims(authPayload)
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

		// Encrypt Claims Using RSA Encryption
		claimsEncrypted, err := hlp.EncryptWithRSA(authClaims["data"].(string))
		if err != nil {
			router.ResponseInternalError(w, err.Error())
			return
		}

		// Set Encrypted Claims to HTTP Header
		r.Header.Set("X-JWT-Claims", claimsEncrypted)

		// Call Next Handler Function With Current Request
		next.ServeHTTP(w, r)
	})
}

// GetJWTToken Function to Generate JWT Token
func GetJWTToken(payload interface{}) (string, error) {
	// Convert Signing Key in Byte Format
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(hlp.KeyRSACfg.BytePrivate)
	if err != nil {
		return "", err
	}

	// Create JWT Token With RS256 Method And Set JWT Claims
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaimsData{
		payload.(string),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	})

	// Generate JWT Token String With Signing Key
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	// Return The JWT Token String and Error
	return tokenString, nil
}

// GetJWTClaims Function to Get JWT Claims in Plain Text
func GetJWTClaims(data string) (string, error) {
	// Decrypt Encrypted Claims Using RSA Encryption
	claimsDecrypted, err := hlp.DecryptWithRSA(data)
	if err != nil {
		return "", err
	}

	// Return Decrypted Claims and Error
	return claimsDecrypted, nil
}

// JWTClaims Function to Get JWT Claims Information
func jwtClaims(data string) (jwt.MapClaims, error) {
	// Convert Signing Key in Byte Format
	signingKey, err := jwt.ParseRSAPublicKeyFromPEM(hlp.KeyRSACfg.BytePublic)
	if err != nil {
		return nil, err
	}

	// Parse JWT Token, If Token is Not Valid Then Return The Signing Key
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	// If Error Found Then Return Empty Claims and The Error
	if err != nil {
		return nil, err
	}

	// Get The Claims
	claims := token.Claims.(jwt.MapClaims)

	// Return The Claims and Error
	return claims, err
}
