package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
)

type CallArgs struct {
	Caller vm.ContractRef
	Addr   common.Address
	Input  []byte
	Gas    uint64
	Value  *big.Int
}

type CallReply struct {
	Ret         []byte
	LeftOverGas uint64
}

type CreateArgs struct {
	Caller vm.ContractRef
	Code   []byte
	Gas    uint64
	Value  *big.Int
}

type CreateReply struct {
	Ret          []byte
	ContractAddr common.Address
	LeftOverGas  uint64
}

// EvmRpcClient is an interface for calling methods on the vm.EVM struct, which
// is used to execute EVM transactions and lives in the SGX binary.
type EvmRpcClient interface {
	// PrepareTx prepares the EVM transaction for execution. It should be the
	// first method called when executing an EVM transaction.
	PrepareTx(args *PrepareTxArgs, reply *PrepareTxReply) error

	// Call is the vm.EVM#Call method.
	Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
	// Call is the vm.EVM#Create method.
	Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
}

// PrepareTxArgs is the argument struct for the statedb.Keeper#PrepareTx method.
type PrepareTxArgs struct {
	// PrepareTxBlockContext is the context of the block in which the transaction
	// will be executed.
	BlockContext PrepareTxBlockContext
	// Msg is the EVM transaction message to run on the EVM.
	Msg core.Message
	// EvmConfig is the EVM configuration to set.
	//
	// IMPORTANT: the tracer field should not be set on this field here,
	// because it is an interface and cannot be passed over RPC.
	EvmConfig EVMConfig
}

type PrepareTxBlockContext struct {
	BlockHeight   int64
	BlockGasLimit uint64
	BlockTime     int64
}

// PrepareTxReply is the reply struct for the statedb.Keeper#PrepareTx method.
type PrepareTxReply struct {
}
