package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/partyhall/partyhall/models"
)

func GetClaimsFromUser(dbUser *models.User) models.JwtCustomClaims {
	return models.JwtCustomClaims{
		Name:     dbUser.Name,
		Username: dbUser.Username,
		Roles:    dbUser.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "PARTYHALL",
			Subject:   fmt.Sprintf("%v", dbUser.Id),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
}

func GetOrGenerateJwtKeys() ([]byte, []byte, error) {
	privateKeyPath := GetPath("partyhall.pem")
	publicKeyPath := GetPath("partyhall.pub")

	hasPrivateKey := true
	hasPublicKey := true

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		hasPrivateKey = false
	}

	if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		hasPublicKey = false
	}

	//#region Generating keys
	if !hasPrivateKey || !hasPublicKey {
		if hasPrivateKey {
			os.Remove(privateKeyPath)
		}

		if hasPublicKey {
			os.Remove(publicKeyPath)
		}

		fmt.Println("Generating JWT keys...")
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, nil, err
		}

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
		err = os.WriteFile(privateKeyPath, privateKeyPEM, 0600)
		if err != nil {
			return nil, nil, err
		}

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return nil, nil, err
		}
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyBytes})
		err = os.WriteFile(publicKeyPath, publicKeyPEM, 0644)
		if err != nil {
			return nil, nil, err
		}
	}
	//#endregion

	//#region Load keys
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}
	//#endregion

	return privateKeyBytes, publicKeyBytes, nil
}
