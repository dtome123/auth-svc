package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// CLI flags
	outDir := flag.String("out-dir", "cert", "Directory to save the RSA key pair")

	flag.Parse()

	if err := generateRSAKeyPair(*outDir); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("RSA key pair generated at:", *outDir)
}

func generateRSAKeyPair(dir string) error {
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	// Generate RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	// Encode private key
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privFilePath := filepath.Join(dir, "private.pem")
	if err := writePEMFile(privFilePath, "RSA PRIVATE KEY", privBytes); err != nil {
		return err
	}

	// Encode public key
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	pubFilePath := filepath.Join(dir, "public.pem")
	if err := writePEMFile(pubFilePath, "PUBLIC KEY", pubBytes); err != nil {
		return err
	}

	return nil
}

func writePEMFile(path, pemType string, bytes []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: pemType, Bytes: bytes}); err != nil {
		return fmt.Errorf("failed to encode %s: %w", path, err)
	}
	return nil
}
