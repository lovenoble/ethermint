package keeper

import (
	"fmt"
	"math/big"
	"net/rpc"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/ethermint/x/evm/statedb"
)

type sgxRpcClient struct {
	logger log.Logger
	cl     *rpc.Client
}

// newSgxRpcClient creates a new RPC client to communicate with the SGX binary.
func newSgxRpcClient(logger log.Logger) (*sgxRpcClient, error) {
	// TODO Make ports configurable
	cl, err := rpc.DialHTTP("tcp", "localhost"+":9092")
	if err != nil {
		return nil, err
	}

	return &sgxRpcClient{
		logger: logger,
		cl:     cl,
	}, nil
}

func (c *sgxRpcClient) makeCall(method string, args, reply any) error {
	c.logger.Debug(fmt.Sprintf("RPC call %s", method), "args", args)
	err := c.cl.Call(method, args, reply)
	if err != nil {
		c.logger.Error(fmt.Sprintf("RPC call %s failed", method), "err", err)
	} else {
		c.logger.Debug(fmt.Sprintf("RPC call %s", method), "reply", reply)
	}
	return err
}

func (c *sgxRpcClient) PrepareTx(args *PrepareTxArgs, reply *PrepareTxReply) error {
	return c.makeCall("SgxRpcServer.PrepareTx", args, reply)
}

func (c *sgxRpcClient) Call(args *CallArgs, reply *CallReply) error {
	return c.makeCall("SgxRpcServer.Call", args, reply)
}

func (c *sgxRpcClient) Create(args *CreateArgs, reply *CreateReply) error {
	return c.makeCall("SgxRpcServer.Create", args, reply)
}

// PrepareTxEVMConfig only contains the fields from EVMConfig that are needed
// to create a new EVM instance. This is used to pass the EVM configuration
// over RPC to the SGX binary.
type PrepareTxEVMConfig struct {
	// ChainConfig is the EVM chain configuration in JSON format. Since the
	// underlying params.ChainConfig struct contains pointer fields, they are
	// not serializable over RPC with gob. Instead, the JSON representation is
	// used.
	ChainConfigJson []byte

	// Fields from EVMConfig
	CoinBase   common.Address
	BaseFee    *big.Int
	TxConfig   statedb.TxConfig
	DebugTrace bool

	// Fields from EVMConfig.FeeMarketParams struct
	NoBaseFee bool

	// Fields from EVMConfig.Params struct
	EvmDenom  string
	ExtraEips []int
}

// PrepareTxArgs is the argument struct for the SgxRpcServer.PrepareTx RPC method.
type PrepareTxArgs struct {
	// Header is the Tendermint header of the block in which the transaction
	// will be executed.
	Header cmtproto.Header
	// Msg is the EVM transaction message to run on the EVM.
	Msg core.Message
	// EvmConfig is the EVM configuration to set.
	EvmConfig PrepareTxEVMConfig
}

// PrepareTxArgs is the reply struct for the SgxRpcServer.PrepareTx RPC method.
type PrepareTxReply struct {
}

// CallArgs is the argument struct for the SgxRpcServer.Call RPC method.
type CallArgs struct {
	Caller vm.AccountRef
	Addr   common.Address
	Input  []byte
	Gas    uint64
	Value  *big.Int
}

// CallReply is the reply struct for the SgxRpcServer.Call RPC method.
type CallReply struct {
	Ret         []byte
	LeftOverGas uint64
}

// CreateArgs is the argument struct for the SgxRpcServer.Create RPC method.
type CreateArgs struct {
	Caller vm.AccountRef
	Code   []byte
	Gas    uint64
	Value  *big.Int
}

// CreateReply is the reply struct for the SgxRpcServer.Create RPC method.
type CreateReply struct {
	Ret          []byte
	ContractAddr common.Address
	LeftOverGas  uint64
}
