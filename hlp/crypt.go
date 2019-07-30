package hlp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
)

// Key RSA Config Struct
type keyRSAConfig struct {
	BytePrivate []byte
	BytePublic  []byte
	KeyPrivate  *rsa.PrivateKey
	KeyPublic   *rsa.PublicKey
}

// KeyRSACfg Variable
var KeyRSACfg keyRSAConfig

// Initialize Function in Helper Cryptography
func init() {
	var err error

	// Load RSA Private Key as Bytes From Private Key File
	KeyRSACfg.BytePrivate, err = ioutil.ReadFile(Config.GetString("CRYPT_PRIVATE_KEY_FILE"))
	if err != nil {
		LogPrintln(LogLevelFatal, "init-crypt", err.Error())
	}

	// Load RSA Private Key Data By Converting RSA Private Key Bytes
	KeyRSACfg.KeyPrivate, err = BytesToPrivateKey(KeyRSACfg.BytePrivate)
	if err != nil {
		LogPrintln(LogLevelFatal, "init-crypt", err.Error())
	}

	// Load RSA Public Key as Bytes From Public Key File
	KeyRSACfg.BytePublic, err = ioutil.ReadFile(Config.GetString("CRYPT_PUBLIC_KEY_FILE"))
	if err != nil {
		LogPrintln(LogLevelFatal, "init-crypt", err.Error())
	}

	// Load RSA Public Key Data By Converting RSA Public Key Bytes
	KeyRSACfg.KeyPublic, err = BytesToPublicKey(KeyRSACfg.BytePublic)
	if err != nil {
		LogPrintln(LogLevelFatal, "init-crypt", err.Error())
	}
}

// BytesToPrivateKey Function
func BytesToPrivateKey(bytePrivate []byte) (*rsa.PrivateKey, error) {
	var err error

	// Decode RSA Private Key Bytes to PEM Block Bytes
	pemBlock, _ := pem.Decode(bytePrivate)
	isEncrypted := x509.IsEncryptedPEMBlock(pemBlock)
	byteBlock := pemBlock.Bytes

	// Check If RSA Private Key Bytes Encrypted
	// If Encrypted Try to Decode it
	if isEncrypted {
		byteBlock, err = x509.DecryptPEMBlock(pemBlock, nil)
		if err != nil {
			return nil, err
		}
	}

	// Parse RSA Private Key Using PKCS1
	rsaKeyPrivate, err := x509.ParsePKCS1PrivateKey(byteBlock)
	if err != nil {
		return nil, err
	}

	// Return RSA Public Key
	return rsaKeyPrivate, nil
}

// BytesToPublicKey Function
func BytesToPublicKey(bytePublic []byte) (*rsa.PublicKey, error) {
	var err error

	// Decode RSA Public Key Bytes to PEM Block Bytes
	pemBlock, _ := pem.Decode(bytePublic)
	isEncrypted := x509.IsEncryptedPEMBlock(pemBlock)
	byteBlock := pemBlock.Bytes

	// Check If RSA Public Key Bytes Encrypted
	// If Encrypted Try to Decode it
	if isEncrypted {
		byteBlock, err = x509.DecryptPEMBlock(pemBlock, nil)
		if err != nil {
			return nil, err
		}
	}

	// Parse RSA Public Key Using PKCIX
	rsaKeyPublic, err := x509.ParsePKIXPublicKey(byteBlock)
	if err != nil {
		return nil, err
	}

	// Return RSA Public Key
	return rsaKeyPublic.(*rsa.PublicKey), nil
}

// EncryptWithRSA Function
func EncryptWithRSA(data string) (string, error) {
	// Generate New SHA512 Hash
	hash := sha512.New()

	// Encrypt Plain Text to Chiper Text Using RSA Encryption OAEP
	chiperText, err := rsa.EncryptOAEP(hash, rand.Reader, KeyRSACfg.KeyPublic, []byte(data), nil)
	if err != nil {
		return "", err
	}

	// Compress Chiper Text to Base64 Format
	compressText := base64.StdEncoding.EncodeToString(chiperText)

	// Return Compressed Text
	return compressText, nil
}

// DecryptWithRSA Function
func DecryptWithRSA(data string) (string, error) {
	// Generate New SHA512 Hash
	hash := sha512.New()

	// Decompress Chiper Text from Base64 Format
	decompressText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Decrypt Chiper Text to Plain Text Using RSA Encryption OAEP
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, KeyRSACfg.KeyPrivate, []byte(decompressText), nil)
	if err != nil {
		return "", err
	}

	// Return Chiper Text
	return string(plainText), nil
}
