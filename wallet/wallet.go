package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/transeptorlabs/betsy/internal/utils"

	"github.com/rs/zerolog/log"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const keyStorePath = "./wallet/tmp"
const coinbaseKeyStorePath = "./wallet/tmp/coinbase/"
const DefaultSeedPhrase = "test test test test test test test test test test test junk"
const cfEntrypointV7Address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"

// BundlerWalletDetails contains the details of the bundler wallet
type BundlerWalletDetails struct {
	Beneficiary       common.Address
	Mnemonic          string
	EntryPointAddress common.Address
}

// Wallet contains the details of the wallet for Besty
type Wallet struct {
	client                    *ethclient.Client
	CoinbaseAddress           common.Address
	BundlerBeneficiaryAddress common.Address
	defaultDevAccounts        []DefaultDevAccount
	keyStore                  *keystore.KeyStore
	password                  string
	EntryPointAddress         common.Address
}

// DefaultDevAccount contains the details of the default development account
type DefaultDevAccount struct {
	Address    common.Address
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

// NewWallet creates a new wallet for Besty
func NewWallet(ctx context.Context, ethNodePort string, coinbaseKeystoreFile string) (*Wallet, error) {
	client, err := ethclient.Dial("http://localhost:" + ethNodePort)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Ethereum client: %v", err)
	}

	// Define the path to the keystore
	ks := keystore.NewKeyStore(keyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	// Load the coinbase account from the keystore
	coinbaseFile := coinbaseKeyStorePath + coinbaseKeystoreFile
	jsonBytes, err := os.ReadFile(coinbaseFile)
	if err != nil {
		return nil, err
	}

	password := ""
	cbAccount, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		return nil, err
	}

	err = utils.RemoveFile(coinbaseFile)
	if err != nil {
		return nil, err
	}

	// Create a new account for the bundler beneficiary
	bAccount, err := createAccount(ks, password)
	if err != nil {
		return nil, err
	}

	defaultDevAccounts, err := GenerateAccountsFromSeed(DefaultSeedPhrase, 10)
	if err != nil {
		return nil, err
	}

	// TODO: Deploy the 4337 pre-compiled contracts

	wallet := &Wallet{
		client:                    client,
		keyStore:                  ks,
		defaultDevAccounts:        defaultDevAccounts,
		CoinbaseAddress:           cbAccount.Address,
		BundlerBeneficiaryAddress: bAccount,
		password:                  password,
		EntryPointAddress:         common.HexToAddress("0x5FbDB2315678afecb367f032d93F642f64180aa3"),
	}

	// Fund the default development accounts
	for _, account := range defaultDevAccounts {
		err = wallet.fundAccountWithEth(ctx, account.Address)
		if err != nil {
			return nil, err
		}
	}

	return wallet, nil
}

// GetKeyStoreAccounts returns all the accounts in the keystore
func (w *Wallet) GetKeyStoreAccounts() []common.Address {
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, w.keyStore)
	return am.Accounts()
}

// PrintDevAccounts prints the default development accounts
func (w *Wallet) PrintDevAccounts(ctx context.Context) error {
	fmt.Println("_________________________Default Development Accounts:_________________________")
	for i, account := range w.defaultDevAccounts {

		balance, err := w.client.BalanceAt(ctx, account.Address, nil)
		if err != nil {
			return err
		}

		fmt.Printf("Account %d:\n", i+1)
		fmt.Printf("Address: %s\n", account.Address.Hex())
		fmt.Printf("Private Key: 0x%s\n", hex.EncodeToString(crypto.FromECDSA(account.PrivateKey)))
		fmt.Printf("Balance: %d wei\n\n", balance)
	}
	fmt.Println("_______________________________________________________________________________")
	return nil
}

// fundAccountWithEth send 4337 ETH to the account using the coinbase account
func (w *Wallet) fundAccountWithEth(ctx context.Context, toAddress common.Address) error {
	// Unlock the account (in the context of the keystore is necessary because the private key is encrypted for security reasons)
	account, err := w.keyStore.Find(accounts.Account{Address: w.CoinbaseAddress})
	if err != nil {
		return err
	}

	err = w.keyStore.Unlock(account, w.password)
	if err != nil {
		return err
	}

	nonce, err := w.client.PendingNonceAt(context.Background(), w.CoinbaseAddress)
	if err != nil {
		return err
	}

	value := new(big.Int)
	value.SetString("4337000000000000000000", 10) // 4337 ETH in Wei
	gasLimit := uint64(21000)                     // The gas limit for a standard ETH transfer is 21000 units.
	gasPrice, err := w.client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction and send it
	chainID, err := w.client.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := w.keyStore.SignTx(account, tx, chainID)
	if err != nil {
		return err
	}

	err = w.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	log.Info().Msgf("tx sent: %s", signedTx.Hash().Hex())

	return nil
}

func (w *Wallet) deploy4337PreCompiledContracts() error {
	return nil
}

// GetAccount generates a new key and stores it into the key directory, encrypting it with the passphrase and return the address
func createAccount(ks *keystore.KeyStore, password string) (common.Address, error) {
	account, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return account.Address, nil
}

// generateAccountsFromSeed generates a number of accounts from a seed phrase
func GenerateAccountsFromSeed(seedPhrase string, numAccounts int) ([]DefaultDevAccount, error) {
	var accounts []DefaultDevAccount

	// Generate the seed from the mnemonic and create a master key from the seed
	seed := bip39.NewSeed(seedPhrase, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive the Ethereum account keys (follows the path m/44'/60'/0'/0/x, where x is the account index)
	// This path is used to generate hierarchical deterministic (HD) wallets.
	for i := 0; i < numAccounts; i++ {
		purpose, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
		if err != nil {
			return nil, err
		}
		coinType, err := purpose.NewChildKey(bip32.FirstHardenedChild + 60)
		if err != nil {
			return nil, err
		}
		account, err := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
		if err != nil {
			return nil, err
		}
		change, err := account.NewChildKey(0)
		if err != nil {
			return nil, err
		}
		key, err := change.NewChildKey(uint32(i))
		if err != nil {
			return nil, err
		}

		// Generate the private key, public key and Ethereum address
		privateKey, err := crypto.ToECDSA(key.Key)
		if err != nil {
			return nil, err
		}

		publicKey := &privateKey.PublicKey
		address := crypto.PubkeyToAddress(*publicKey)
		devAccount := DefaultDevAccount{
			Address:    address,
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}

		accounts = append(accounts, devAccount)
	}

	return accounts, nil
}
