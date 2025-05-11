package controller

import (
	"crypto/ecdsa"
	"encoding/hex"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateWallet() (string, string) {
	generatedPrivateKey, err := crypto.GenerateKey()

	if err != nil {
		log.Fatalf("Error while generating private %v", err)
	}

	publicKey := generatedPrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		log.Fatal("Error casting public key to ECDSA")
	}

	ethAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	privateKeyBytes := crypto.FromECDSA(generatedPrivateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	return privateKeyHex, ethAddress
}
