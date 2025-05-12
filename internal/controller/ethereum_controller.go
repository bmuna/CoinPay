package controller

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/bmuna/CoinPay/backend/internal/models"
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

func SendEth(w http.ResponseWriter, r *http.Request) {
	infuraURL := os.Getenv("INFURA_URL")
	privateKeyHex := os.Getenv("PRIVATE_KEY")

	// Parse JSON request body
	var req models.SendEthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Connect to Ethereum node
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Convert private key hex to ECDSA
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Derive public key and sender address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Prepare recipient address and amount
	toAddress := common.HexToAddress(req.ToAddress)
	amountInWei := new(big.Int)
	amountInWei.Mul(big.NewInt(int64(req.Amount*1e6)), big.NewInt(1e12))

	ctx := context.Background()

	// Check balance (optional)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	if balance.Cmp(amountInWei) < 0 {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Get nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Create transaction
	gasLimit := uint64(21000) // Standard for ETH transfer
	tx := types.NewTransaction(nonce, toAddress, amountInWei, gasLimit, gasPrice, nil)

	// Sign transaction
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	// Send transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	// Return TX hash
	fmt.Fprintf(w, "Transaction sent! TX Hash: %s", signedTx.Hash().Hex())
}
