package transition

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func ApplyAccountToState(state *trie.Database, acckey string, acc *types.StateAccount) error {
	// We are using regular trie sice StateTrie provides a wrapper around a Trie that hashes the addrss of the accounts to make them node keys, while we already have hashes only
	accTrie, err := trie.New(trie.StateTrieID(types.EmptyRootHash), state)
	if err != nil {
		return err
	}
	val, _ := rlp.EncodeToBytes(acc)
	err = accTrie.Update([]byte(acckey), val)
	if err != nil {
		return err
	}
	return nil
}
