package keeper

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// EvmRpcClient is an interface for calling methods on the vm.EVM struct, which
// is used to execute EVM transactions and lives in the SGX binary.
type EvmRpcClient interface {
	// Call is the vm.EVM#Call method.
	Call(caller vm.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
	// Call is the vm.EVM#Create method.
	Create(caller vm.ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
}
