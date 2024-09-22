package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"strings"
)

const (
	PRIVPKCS1 = "RSA PRIVATE KEY"
	PRIVPKCS8 = "PRIVATE KEY"

	PUBPKCS1    = "RSA PUBLIC KEY"
	PUBPKCS8    = "PUBLIC KEY"
	CERTIFICATE = "CERTIFICATE"
)

type (
	IRsaCrypto interface {
		PrivateKey(value string) error
		PublicKey(value string, raw bool) ([]byte, error)
		Asymmetric(req *Asymmetric, password string) error
	}

	Asymmetric struct {
		ClientID   string `json:"clientId"`
		ClientKey  string `json:"clientKey"`
		PrivateKey string `json:"privateKey"`
	}

	rsaCrypto struct{}
)

func NewRsa() IRsaCrypto {
	return &rsaCrypto{}
}

func (h *rsaCrypto) PrivateKey(value string) error {
	var privateKey string

	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}

	pemDecoded, _ := pem.Decode([]byte(decoded))
	if pemDecoded == nil {
		return errors.New("Invalid PEM PrivateKey certificate")
	}

	if pemDecoded.Type == PRIVPKCS1 {
		privateKey = string(pem.EncodeToMemory(pemDecoded))
	} else if pemDecoded.Type == PRIVPKCS8 {
		privateKey = string(pem.EncodeToMemory(pemDecoded))
	} else if pemDecoded.Type == CERTIFICATE {
		privateKey = string(pem.EncodeToMemory(pemDecoded))
	} else {
		return errors.New("Invalid PEM PrivateKey certificate")
	}

	if privateKey == "" {
		return errors.New("Invalid PEM PrivateKey certificate")
	}

	return nil
}

func (h *rsaCrypto) PublicKey(value string, rawPem bool) ([]byte, error) {
	var publicKey []byte

	externalPublicKey, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return nil, err
	}

	pemDecoded, _ := pem.Decode([]byte(externalPublicKey))
	if pemDecoded == nil {
		return nil, errors.New("Invalid PEM PublicKey certificate")
	}

	if !rawPem && pemDecoded.Type == PUBPKCS1 {
		publicKey = pem.EncodeToMemory(pemDecoded)
	} else if !rawPem && pemDecoded.Type == PUBPKCS8 {
		publicKey = pem.EncodeToMemory(pemDecoded)
	} else if !rawPem && pemDecoded.Type == CERTIFICATE {
		publicKey = pem.EncodeToMemory(pemDecoded)
	} else {
		publicKey = pemDecoded.Bytes
	}

	return publicKey, nil
}

func (h *rsaCrypto) Asymmetric(req *Asymmetric, password string) error {
	salt := rand.Reader
	rsaPrivateKey := new(rsa.PrivateKey)

	internalClientID := os.Getenv("INTERNAL_CLIENT_ID")
	internalClientKey := os.Getenv("INTERNAL_CLIENT_KEY")
	internalPrivateKey := os.Getenv("INTERNAL_PRIVATE_KEY")
	internalPublicKey := os.Getenv("INTERNAL_PUBLIC_KEY")

	headers := Asymmetric{}
	headers.ClientID = strings.TrimSpace(req.ClientID)
	headers.ClientKey = strings.TrimSpace(req.ClientKey)
	headers.PrivateKey = strings.TrimSpace(req.PrivateKey)

	cipherBody := []byte(headers.ClientID + ":" + headers.ClientKey + ":" + headers.PrivateKey)
	cipherBodyHash256 := sha256.New()
	cipherBodyHash256.Write(cipherBody)
	cipherBodyHash := cipherBodyHash256.Sum(nil)

	if headers.PrivateKey == "" {
		return errors.New("PEM PrivateKey certificate must be a base64 format")
	} else if headers.ClientID != internalClientID {
		return errors.New("Invalid ClientId")
	} else if headers.ClientKey != internalClientKey {
		return errors.New("Invalid ClientKey")
	} else if strings.Compare(headers.PrivateKey, internalPrivateKey) != 0 {
		return errors.New("Invalid PrivateKey")
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(headers.PrivateKey)
	if err != nil {
		return err
	}

	headers.PrivateKey = string(decodedPrivateKey)
	pemDecode, _ := pem.Decode([]byte(headers.PrivateKey))

	if pemDecode == nil {
		return errors.New("Invalid PEM PrivateKey certificate")
	}

	switch pemDecode.Type {
	case PRIVPKCS1:
		if password != "" && x509.IsEncryptedPEMBlock(pemDecode) {
			pemBlockDecrypt, err := x509.DecryptPEMBlock(pemDecode, []byte(strings.TrimSpace(password)))
			if err != nil {
				return err
			}

			parsePrivateKey, err := x509.ParsePKCS1PrivateKey(pemBlockDecrypt)
			if err != nil {
				return err
			}

			rsaPrivateKey = parsePrivateKey
		}

		parsePrivateKey, err := x509.ParsePKCS1PrivateKey(pemDecode.Bytes)
		if err != nil {
			return err
		}

		rsaPrivateKey = parsePrivateKey

		break

	case PRIVPKCS8:
		if password != "" && x509.IsEncryptedPEMBlock(pemDecode) {
			pemBlockDecrypt, err := x509.DecryptPEMBlock(pemDecode, []byte(strings.TrimSpace(password)))
			if err != nil {
				return err
			}

			parsePrivateKey, err := x509.ParsePKCS8PrivateKey(pemBlockDecrypt)
			if err != nil {
				return err
			}

			rsaPrivateKey = parsePrivateKey.(*rsa.PrivateKey)
		}

		parsePrivateKey, err := x509.ParsePKCS8PrivateKey(pemDecode.Bytes)
		if err != nil {
			return err
		}

		rsaPrivateKey = parsePrivateKey.(*rsa.PrivateKey)

		break

	default:
		return errors.New("Invalid PEM PrivateKey certificate")
	}

	if err := rsaPrivateKey.Validate(); err != nil {
		return err
	}

	signature, err := rsa.SignPKCS1v15(salt, rsaPrivateKey, crypto.SHA256, cipherBodyHash)
	if err != nil {
		return err
	}

	rsaPublicKeyRaw, err := h.PublicKey(internalPublicKey, true)
	if err != nil {
		return err
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(rsaPublicKeyRaw)
	if err != nil {
		return err
	}

	if err := rsa.VerifyPKCS1v15(rsaPublicKey.(*rsa.PublicKey), crypto.SHA256, cipherBodyHash, signature); err != nil {
		return err
	}

	return nil
}
