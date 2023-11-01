package transition

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	parquet "github.com/parquet-go/parquet-go"
	"io"
	"os"
)

func LoadTransactionsFromParquet(src string) (types.Transactions, error) {
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}

	schema := parquet.SchemaOf(new(parquetTxRow))
	reader := parquet.NewReader(f, schema)
	defer reader.Close()

	txs := make([]*types.Transaction, 0)
	row := new(parquetTxRow)

	for {
		err = reader.Read(row)
		if err != nil {
			if err == io.EOF {
				log.Debug("EOF")
				break
			}
			return nil, fmt.Errorf("error reading parquet file form source %s with err: %w", src, err)
		}

		tx := &types.Transaction{}
		err = tx.UnmarshalBinary(row.RawTx)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling rawTx to Transction type with err: %w", err)
		}
		/// We can add filtering logic here
		txs = append(txs, tx)
	}
	return txs, nil
}

type parquetTxRow struct {
	Timestamp              uint     `parquet:"timestamp"`
	Hash                   string   `parquet:"hash"`
	ChanId                 string   `parquet:"chanId"`
	From                   string   `parquet:"from"`
	To                     string   `parquet:"to"`
	Value                  string   `parquet:"value"`
	Nonce                  string   `parquet:"nonce"`
	Gas                    string   `parquet:"gas"`
	GasPrice               string   `parquet:"gasPrice"`
	GasTipCap              string   `parquet:"gasTipCap"`
	GasFeeCap              string   `parquet:"gasFeeCap"`
	DataSize               int      `parquet:"dataSize"`
	Data4Bytes             string   `parquet:"data4Bytes"`
	Sources                []string `parquet:"sources"`
	IncludeAtBlockHeight   int      `parquet:"includeAtBlockHeight"`
	IncludedBlockTimestamp int      `parquet:"includedBlockTimestamp"`
	InclusionDelayMs       int      `parquet:"inclusionDelayMs"`
	RawTx                  []byte   `parquet:"rawTx"`
}
