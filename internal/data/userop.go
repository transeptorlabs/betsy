package data

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/transeptorlabs/betsy/contracts/entrypoint"
	"github.com/transeptorlabs/betsy/internal/utils"
)

// UserOpV7Hexify is a struct used to store user operations in a database. Contains all EIP-4337 with all hex fields.
type UserOpV7Hexify struct {
	Sender string `json:"sender"               mapstructure:"sender"               validate:"required"`
	Nonce  string `json:"nonce"               mapstructure:"nonce"               validate:"required"`

	// (optional)
	Factory     string `json:"factory"               mapstructure:"factory" validate:"required"`
	FactoryData string `json:"factoryData"               mapstructure:"factoryData"`

	CallData     string `json:"callData"               mapstructure:"callData" validate:"required"`
	CallGasLimit string `json:"callGasLimit"               mapstructure:"callGasLimit" validate:"required"`

	VerificationGasLimit string `json:"verificationGasLimit"               mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   string `json:"preVerificationGas"               mapstructure:"preVerificationGas" validate:"required"`

	MaxFeePerGas         string `json:"maxFeePerGas"               mapstructure:"maxFeePerGas" validate:"required" `
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"               mapstructure:"maxPriorityFeePerGas" validate:"required" `

	// (optional)
	Paymaster                     string `json:"paymaster"     mapstructure:"paymaster"     validate:"required"`
	PaymasterVerificationGasLimit string `json:"paymasterVerificationGasLimit"     mapstructure:"paymasterVerificationGasLimit"`
	PaymasterPostOpGasLimit       string `json:"paymasterPostOpGasLimit"     mapstructure:"paymasterPostOpGasLimit"`
	PaymasterData                 string `json:"paymasterAndData"     mapstructure:"paymasterAndData"`

	Signature string `json:"signature"               mapstructure:"signature"`
}

func (op *UserOpV7Hexify) GetUserOpHash(epAddress common.Address, ethClient *ethclient.Client) (common.Hash, error) {
	ep, err := entrypoint.NewEntryPointV7(epAddress, ethClient)
	if err != nil {
		return common.Hash{}, err
	}

	packedOp, err := op.packUserOp()
	if err != nil {
		return common.Hash{}, err
	}

	hash, err := ep.GetUserOpHash(&bind.CallOpts{}, *packedOp)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(hash[:]), nil
}

func (op *UserOpV7Hexify) getInitCode() ([]byte, error) {
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

func (op *UserOpV7Hexify) getAccountGasLimits() ([32]byte, error) {
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

func (op *UserOpV7Hexify) gasFees() ([32]byte, error) {
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

func (op *UserOpV7Hexify) getPaymasterAndData() ([]byte, error) {
	if op.Paymaster == "0x" {
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

func (op *UserOpV7Hexify) packUserOp() (*entrypoint.PackedUserOperation, error) {
	nonceDecoded, err := hexutil.DecodeBig(op.Nonce)
	if err != nil {
		return nil, errors.New("nonce (bigInt) conversion failed")
	}

	initCodeDecoded, err := op.getInitCode()
	if err != nil {
		return nil, errors.New("nonce (bigInt) conversion failed")
	}

	callDataDecoded, err := hexutil.Decode(op.CallData)
	if err != nil {
		return nil, errors.New("calldata (bytes) conversion failed")
	}

	preVerificationGaseDecoded, err := hexutil.DecodeBig(op.PreVerificationGas)
	if err != nil {
		return nil, errors.New("preVerificationGas (bigInt) conversion failed")
	}

	accountGasLimitsDecoded, err := op.getAccountGasLimits()
	if err != nil {
		return nil, errors.New("accountGasLimit (bytes 32) conversion failed")
	}

	gasFeesDecoded, err := op.gasFees()
	if err != nil {
		return nil, errors.New("gasFees (bytes 32) conversion failed")
	}

	paymasterAndDataDecoded, err := op.getPaymasterAndData()
	if err != nil {
		return nil, errors.New("paymasterAndData (bytes) conversion failed")
	}

	signatureDecoded, err := hexutil.Decode(op.Signature)
	if err != nil {
		return nil, errors.New("signature (bytes) conversion failed")
	}

	return &entrypoint.PackedUserOperation{
		Sender:             common.HexToAddress(op.Sender),
		Nonce:              nonceDecoded,
		InitCode:           initCodeDecoded,
		CallData:           callDataDecoded,
		AccountGasLimits:   accountGasLimitsDecoded,
		PreVerificationGas: preVerificationGaseDecoded,
		GasFees:            gasFeesDecoded,
		PaymasterAndData:   paymasterAndDataDecoded,
		Signature:          signatureDecoded,
	}, nil
}
