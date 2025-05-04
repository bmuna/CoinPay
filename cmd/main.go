package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/bmuna/CoinPay/backend/internal/eth"
	"github.com/bmuna/CoinPay/backend/internal/server"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	privateKeyHex := os.Getenv("PRIVATE_KEY")
	infuraURL := os.Getenv("INFURA_URL")
	usdtAddress := os.Getenv("USDT_ADDRESS")
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port number is empty")
	}

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Convert your private key to an *ecdsa.PrivateKey
	privateKey, err := crypto.HexToECDSA(privateKeyHex)

	// fromAddress := common.HexToAddress("0xdBe3098CCBF320256b62F948B5400bEB6661be96")
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Replace with the actual address you want to send to
	toAddress := common.HexToAddress("0xReceiverAddressHere")

	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// Send ETH (replace the amount with how much ETH you want to send, in wei)
	amount := "1000000000000000000" // 1 ETH in wei

	// Call the sendETH function to send the transaction
	eth.SendETH(client, privateKey, fromAddress, toAddress, amount)

	// Optionally, get balances (ETH or ERC-20 token balance) after sending
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("Failed to retrieve the balance: %v", err)
	}

	etherBalance := new(big.Float).SetInt(balance)
	etherBalance.Quo(etherBalance, big.NewFloat(1e18))

	fmt.Printf("Balance of ETH for address %s: %.6f ETH\n", fromAddress.Hex(), etherBalance)

	// Get USDT balance using the ERC-20 token contract
	// usdtAddress := " // USDT contract on Ethereum mainnet
	tokenBalance, err := eth.GetERC20Balance(client, usdtAddress, fromAddress.Hex())
	if err != nil {
		log.Fatalf("Failed to retrieve USDT balance: %v", err)
	}

	fmt.Printf("Balance of USDT for address %s: %.6f USDT\n", fromAddress.Hex(), tokenBalance)

	handler := server.NewServer()

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
