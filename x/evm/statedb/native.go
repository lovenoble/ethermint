package statedb

import (
	"cosmossdk.io/store/types"
	"github.com/ethereum/go-ethereum/common"
)

var _ JournalEntry = nativeChange{}

type nativeChange struct {
	snapshot types.MultiStore
	events   int
}

func (native nativeChange) Dirtied() *common.Address {
	return nil
}

func (native nativeChange) Revert(s *StateDB) {
	s.restoreNativeState(native.snapshot)
	s.nativeEvents = s.nativeEvents[:len(s.nativeEvents)-native.events]
}
