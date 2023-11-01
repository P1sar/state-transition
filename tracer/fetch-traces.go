package tracer

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"math/big"
)

func FetchTracesForBlock(url string, block *big.Int) []BlockTraces {

	result := make([]BlockTraces, 0)

	rpcClient, err := rpc.DialContext(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}

	if err = rpcClient.CallContext(context.Background(), &result, "trace_block", hexutil.EncodeBig(block)); err != nil {
		log.Fatal(err)
	}
	return result
}

type BlockTraces struct {
	Action              Action         `json:"action,omitempty"`
	BlockHash           string         `json:"blockHash,omitempty"`
	BlockNumber         int            `json:"blockNumber,omitempty"`
	Result              Result         `json:"result,omitempty"`
	Address             common.Address `json:"address,omitempty"`
	Subtraces           int            `json:"subtraces,omitempty"`
	TraceAddress        []any          `json:"traceAddress,omitempty"`
	TransactionHash     string         `json:"transactionHash,omitempty"`
	TransactionPosition int            `json:"transactionPosition,omitempty"`
	Type                string         `json:"type,omitempty"`
}

type Action struct {
	CallType      string         `json:"callType,omitempty"`
	From          common.Address `json:"from,omitempty"`
	Gas           string         `json:"gas,omitempty"`
	Input         string         `json:"input,omitempty"`
	To            common.Address `json:"to,omitempty"`
	Value         string         `json:"value,omitempty"`
	Init          string         `json:"init,omitempty"`
	RefundAddress common.Address `json:"refundAddress,omitempty"`
	Address       common.Address `json:"address,omitempty"`
	Author        common.Address `json:"author,omitempty"`
}

type Result struct {
	Address common.Address `json:"address,omitempty"`
	GasUsed string         `json:"gasUsed,omitempty"`
	Output  string         `json:"output,omitempty"`
}

func ReturnTouchedAddresses(traces []BlockTraces) []common.Address {
	addresseMap := make(map[common.Address]bool)
	for _, t := range traces {
		switch t.Type {
		case "reward":
			addresseMap[t.Action.Author] = true
		case "call":
			addresseMap[t.Action.From] = true
			addresseMap[t.Action.To] = true
		case "create":
			addresseMap[t.Action.From] = true
			addresseMap[t.Result.Address] = true
		case "suicide":
			addresseMap[t.Action.RefundAddress] = true
			addresseMap[t.Action.Address] = true
		}
	}
	addresses := make([]common.Address, 0)
	for k, _ := range addresseMap {
		addresses = append(addresses, k)
	}
	log.Printf("Touched addresss %s", len(addresses))
	return addresses
}
