package crypto_eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
)

func SendETH(client *ethclient.Client, privateKey *ecdsa.PrivateKey, fromAddress, toAddress common.Address, amount string) {
	// Convert amount from string to *big.Int (in wei)
	weiAmount := new(big.Int)
	weiAmount, ok := weiAmount.SetString(amount, 10)
	if !ok {
		log.Fatalf("Failed to convert amount to wei")
	}

	// Fetch the nonce for the sender address
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to retrieve nonce: %v", err)
	}

	// Set the gas price (network fee)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to get gas price: %v", err)
	}

	// Prepare the transaction
	tx := types.NewTransaction(
		nonce,
		toAddress,
		weiAmount,
		21000, // Standard gas limit for a simple ETH transfer
		gasPrice,
		nil, // Data field, set to nil for ETH transfer
	)

	// Sign the transaction using the private key
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("Failed to sign the transaction: %v", err)
	}

	// Send the signed transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction successfully sent: %s\n", signedTx.Hash().Hex())
}
