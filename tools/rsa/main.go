package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {

	os.Mkdir("cert", 0755)

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privFile, _ := os.Create("cert/private.pem")
	pem.Encode(privFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	privFile.Close()

	pubBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	pubFile, _ := os.Create("cert/public.pem")
	pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	pubFile.Close()
}
