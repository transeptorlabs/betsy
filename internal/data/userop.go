package data

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/transeptorlabs/betsy/contracts/entrypoint"
	"github.com/transeptorlabs/betsy/internal/utils"
)

// UserOpV7Hexify is a struct used to store user operations in a database. Contains all EIP-4337 with all hex fields.
type UserOpV7Hexify struct {
	Sender string `json:"sender"               mapstructure:"sender"`
	Nonce  string `json:"nonce"               mapstructure:"nonce"`

	// (optional)
	Factory     string `json:"factory"               mapstructure:"factory"`
	FactoryData string `json:"factoryData"               mapstructure:"factoryData"`

	CallData     string `json:"callData"               mapstructure:"callData"`
	CallGasLimit string `json:"callGasLimit"               mapstructure:"callGasLimit"`

	VerificationGasLimit string `json:"verificationGasLimit"               mapstructure:"verificationGasLimit"`
	PreVerificationGas   string `json:"preVerificationGas"               mapstructure:"preVerificationGas"`

	MaxFeePerGas         string `json:"maxFeePerGas"               mapstructure:"maxFeePerGas" `
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"               mapstructure:"maxPriorityFeePerGas" `

	// (optional)
	Paymaster                     string `json:"paymaster"     mapstructure:"paymaster"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"     mapstructure:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"     mapstructure:"paymasterPostOpGasLimit"`
	PaymasterData                 string `json:"paymasterAndData"     mapstructure:"paymasterAndData"`

	Signature string `json:"signature"               mapstructure:"signature"`
}

// GetUserOpHash returns the hash of the user operation
func (op *UserOpV7Hexify) GetUserOpHash(epAddress common.Address, ethClient *ethclient.Client) (common.Hash, error) {
	ep, err := entrypoint.NewEntryPointV7(epAddress, ethClient)
	if err != nil {
		return common.Hash{}, err
	}

	packedOp, err := op.PackUserOp()
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := ep.GetUserOpHash(&bind.CallOpts{}, *packedOp)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(hash[:]), nil
}

// GetInitCode returns the init code of the user operation
func (op *UserOpV7Hexify) GetInitCode() ([]byte, error) {
	if op.Factory == "0x" {
		return []byte{}, nil
	}

	if op.FactoryData == "0x" {
		return nil, errors.New("Got Factory but missing FactoryData")
	}

	factoryDecoded, err := hexutil.Decode(op.Factory)
	if err != nil {
		return nil, errors.New("factory (bytes) conversion failed")
	}

	factoryDataDecoded, err := hexutil.Decode(op.FactoryData)
	if err != nil {
		return nil, errors.New("factoryData (bytes) conversion failed")
	}

	concatenatedBytes := append(factoryDecoded, factoryDataDecoded...)

	return concatenatedBytes, nil
}

// GetAccountGasLimits returns the account gas limits of the user operation
func (op *UserOpV7Hexify) GetAccountGasLimits() ([32]byte, error) {
	verificationGasLimitDecoded, err := hexutil.Decode(op.VerificationGasLimit)
	if err != nil {
		return [32]byte{}, errors.New("verificationGasLimit (bytes) conversion failed")
	}

	callGasLimitDecoded, err := hexutil.Decode(op.CallGasLimit)
	if err != nil {
		return [32]byte{}, errors.New("callGasLimit (bytes) conversion failed")
	}

	// Truncate if longer than 16 bytes or Pad with leading zeros if shorter than 16 bytes
	verificationGasLimitDecodedPadded, err := utils.PadToBytes16(verificationGasLimitDecoded)
	if err != nil {
		return [32]byte{}, err
	}

	callGasLimitDecodedPadded, err := utils.PadToBytes16(callGasLimitDecoded)
	if err != nil {
		return [32]byte{}, err
	}

	// Concatenate the byte slices
	concatenatedBytes := append(verificationGasLimitDecodedPadded, callGasLimitDecodedPadded...)

	if len(concatenatedBytes) != 32 {
		return [32]byte{}, errors.New("concatenatedBytes(verificationGasLimitDecodedPadded, callGasLimitDecodedPadded) is not equal to 32 bytes")
	}

	// Convert concatenatedBytes to [32]byte
	var result [32]byte
	copy(result[:], concatenatedBytes)

	return result, nil
}

// GasFees returns the gas fees of the user operation
func (op *UserOpV7Hexify) GasFees() ([32]byte, error) {
	maxPriorityFeePerGasDecoded, err := hexutil.Decode(op.MaxPriorityFeePerGas)
	if err != nil {
		return [32]byte{}, errors.New("maxPriorityFeePerGas (bytes) conversion failed")
	}

	maxFeePerGasDecoded, err := hexutil.Decode(op.MaxFeePerGas)
	if err != nil {
		return [32]byte{}, errors.New("maxFeePerGas (bytes) conversion failed")
	}

	// Truncate if longer than 16 bytes or Pad with leading zeros if shorter than 16 bytes
	maxPriorityFeePerGasDecodedPadded, err := utils.PadToBytes16(maxPriorityFeePerGasDecoded)
	if err != nil {
		return [32]byte{}, err
	}

	maxFeePerGasDecodedPadded, err := utils.PadToBytes16(maxFeePerGasDecoded)
	if err != nil {
		return [32]byte{}, err
	}

	// Concatenate the byte slices
	concatenatedBytes := append(maxPriorityFeePerGasDecodedPadded, maxFeePerGasDecodedPadded...)

	if len(concatenatedBytes) != 32 {
		return [32]byte{}, errors.New("concatenatedBytes(maxPriorityFeePerGasDecodedPadded, maxFeePerGasDecodedPadded) is not equal to 32 bytes")
	}

	// Convert concatenatedBytes to [32]byte
	var result [32]byte
	copy(result[:], concatenatedBytes)

	return result, nil
}

// GetPaymasterAndData returns the paymaster and data of the user operation
func (op *UserOpV7Hexify) GetPaymasterAndData() ([]byte, error) {
	if op.Paymaster == "0x" || op.Paymaster == "" {
		return []byte{}, nil
	}

	if op.PreVerificationGas == "0x" {
		return nil, errors.New("Got Paymaster but missing PreVerificationGas")
	}

	if op.PaymasterPostOpGasLimit == "0x" {
		return nil, errors.New("Got Paymaster but missing PaymasterPostOpGasLimit")
	}

	// Decode all paymaster values
	paymasterDecoded, err := hexutil.Decode(op.Paymaster)
	if err != nil {
		return nil, errors.New("paymasterDecoded (bytes) conversion failed")
	}

	preVerificationGasDecoded, err := hexutil.Decode(op.PreVerificationGas)
	if err != nil {
		return []byte{}, errors.New("preVerificationGas (bytes) conversion failed")
	}

	paymasterPostOpGasLimitDecoded, err := hexutil.Decode(op.PaymasterPostOpGasLimit)
	if err != nil {
		return []byte{}, errors.New("paymasterPostOpGasLimit (bytes) conversion failed")
	}

	// Truncate if longer than 16 bytes or Pad with leading zeros if shorter than 16 bytes
	preVerificationGasDecodedPadded, err := utils.PadToBytes16(preVerificationGasDecoded)
	if err != nil {
		return []byte{}, err
	}

	paymasterPostOpGasLimitDecodedPadded, err := utils.PadToBytes16(paymasterPostOpGasLimitDecoded)
	if err != nil {
		return []byte{}, err
	}
	concatenatedGasBytes := append(preVerificationGasDecodedPadded, paymasterPostOpGasLimitDecodedPadded...)

	// Concatenate the byte slices
	concatenatedPaymasterBytes := append(
		paymasterDecoded,
		concatenatedGasBytes...,
	)

	var paymasterDataDecoded []byte
	if op.PaymasterData == "0x" {
		paymasterDataDecoded = []byte{}
	} else {
		paymasterDataDecoded, err = hexutil.Decode(op.PaymasterData)
		if err != nil {
			return nil, errors.New("paymasterData (bytes) conversion failed")
		}
	}

	return append(
		concatenatedPaymasterBytes,
		paymasterDataDecoded...,
	), nil
}

// PackUserOp packs the user operation into a PackedUserOperation struct
func (op *UserOpV7Hexify) PackUserOp() (*entrypoint.PackedUserOperation, error) {
	nonceDecoded := new(big.Int)
	nonceDecoded.SetString(op.Nonce, 16)

	initCodeDecoded, err := op.GetInitCode()
	if err != nil {
		return nil, fmt.Errorf("initCode conversion failed: %w", err)
	}

	callDataDecoded, err := hexutil.Decode(op.CallData)
	if err != nil {
		return nil, fmt.Errorf("calldata (bytes) conversion failed: %w", err)
	}

	preVerificationGasDecoded, err := hexutil.DecodeBig(op.PreVerificationGas)
	if err != nil {
		return nil, fmt.Errorf("preVerificationGas (bigInt) conversion failed: %w", err)
	}

	accountGasLimitsDecoded, err := op.GetAccountGasLimits()
	if err != nil {
		return nil, fmt.Errorf("accountGasLimit (bytes 32) conversion failed: %w", err)
	}

	gasFeesDecoded, err := op.GasFees()
	if err != nil {
		return nil, fmt.Errorf("gasFees (bytes 32) conversion failed: %w", err)
	}

	paymasterAndDataDecoded, err := op.GetPaymasterAndData()
	if err != nil {
		return nil, fmt.Errorf("paymasterAndData (bytes) conversion failed: %w", err)
	}

	signatureDecoded, err := hexutil.Decode(op.Signature)
	if err != nil {
		return nil, fmt.Errorf("signature (bytes) conversion failed: %w", err)
	}

	return &entrypoint.PackedUserOperation{
		Sender:             common.HexToAddress(op.Sender),
		Nonce:              nonceDecoded,
		InitCode:           initCodeDecoded,
		CallData:           callDataDecoded,
		AccountGasLimits:   accountGasLimitsDecoded,
		PreVerificationGas: preVerificationGasDecoded,
		GasFees:            gasFeesDecoded,
		PaymasterAndData:   paymasterAndDataDecoded,
		Signature:          signatureDecoded,
	}, nil
}
