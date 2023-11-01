package transition

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"log"
	"math/big"
	"os"
)

func ReadSnapshotData() *trie.Database {
	file, err := os.Open("/Users/pisar/Downloads/snapshot-dump-10-5.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lvlDB, err := rawdb.NewLevelDBDatabase("/var/www/test-assigments/stateDB", 2048, 1024, "", false)
	defer lvlDB.Close()
	if err != nil {
		log.Fatal(err)
	}
	trieDB := trie.NewDatabase(lvlDB, &trie.Config{Preimages: false})
	defer trieDB.Close()

	accTrie, err := trie.New(trie.StateTrieID(types.EmptyRootHash), trieDB)

	//statedb.DumpToCollector()
	scanner := bufio.NewScanner(file)
	maxBufferSize := 1024 * 1024 * 2048
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	for i := 0; i < 100; i += 1 {
		scanner.Scan()
		fmt.Printf("%s", scanner.Bytes())
		dumpAccount := new(state.DumpAccount)
		err := json.Unmarshal(scanner.Bytes(), dumpAccount)
		if err != nil {
			//fmt.Printf("skpping the line to unmarshal : %s", scanner.Bytes())
			continue
		}
		if dumpAccount.SecureKey.String() != "" {
			balance := big.NewInt(0)
			balance.SetString(dumpAccount.Balance, 10)

			acc := &types.StateAccount{
				Nonce:    dumpAccount.Nonce,
				Balance:  balance,
				Root:     common.BytesToHash(dumpAccount.Root),
				CodeHash: dumpAccount.CodeHash,
			}
			val, _ := rlp.EncodeToBytes(acc)
			err = accTrie.Update(dumpAccount.SecureKey, val)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("account key is nill %s", scanner.Bytes())
			continue
		}
		//fmt.Printf("account applied %s \n", dumpAccount.SecureKey.String())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("Sucesfully over")
	root, nodes, err := accTrie.Commit(true)
	if err != nil {
		log.Fatal(err)
	}
	err = trieDB.Commit(root, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(root)
	log.Println(len(nodes.Leaves))

	return nil
}

func OpenTrieDB() {
	lvlDB, err := rawdb.NewLevelDBDatabase("/var/www/test-assigments/stateDB", 2048, 1024, "", true)
	if err != nil {
		log.Fatal(err)
	}
	trieDB := trie.NewDatabase(lvlDB, &trie.Config{Preimages: true})

	//trieDB.

	accTrie, err := trie.New(trie.StateTrieID(types.EmptyRootHash), trieDB)
	if err != nil {
		log.Fatal(err)
	}
	v, err := accTrie.Get(common.Hex2Bytes("0x000006f88dbc1fc2eb8c072981a032e864f6e05d72176c0f77a01ab6f402187b"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%x", v)
	dumpAccount := new(state.DumpAccount)
	err = json.Unmarshal(v, dumpAccount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dumpAccount.Nonce)
}

// Need to understand what means ` leveraging t8n ` bcs that tasks says produce post state, while Wasif told that it is about testing outcomes of the algorithms that are black boxes
// FUnction should accept list of transctions + old state and then validate them and assign to extend where they give the most profit to the block builder
//   block=18,286,186 hash=0xb22552a5fd7c66e71756585b292ccd51ef43b57c9e614a3db108b8cedfdf8483
//  skipcode=false skipstorage=false start=0x0000000000000000000000000000000000000000000000000000000000000000 limit=0skpping the line to unmarshal : INFO [10-05|21:51:32.474] Snapshot dumping started
//  root=d322d2..8e1dddskpping

//0xea038c191b34bb4ba3342430e0e33501d7176b5db7341e7643fd32ea861f41c6
//{"balance":"162547888961104","nonce":2,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470","key":"0x000006f88dbc1fc2eb8c072981a032e864f6e05d72176c0f77a01ab6f402187b"}
