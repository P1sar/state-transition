package tracer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"errors"
	"fmt"
	"log"
)

var (
	rootCMD = &cobra.Command{
		Use:   "",
		Short: "Fetch provided block traces",
		Long:  "Calls ",
		RunE:  Run,
	}
)

func init() {
	rootCMD.PersistentFlags().String("block", "", "block number")
	rootCMD.PersistentFlags().String("rpc", "", "URL of rpc to fetch the traces")
	rootCMD.PersistentFlags().String("snapshot", "", "Path to the snapshot")

	rootCMD.MarkPersistentFlagRequired("block")
	rootCMD.MarkPersistentFlagRequired("rpc")
	//rootCMD.MarkPersistentFlagRequired("snapshot")

}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Fatalf("failed to execute root cmd with err: %s", err)
	}
}

func Run(cmd *cobra.Command, args []string) error {
	block, err := cmd.Flags().GetString("block")
	if err != nil {
		return err
	}
	blockInt := validateBlockFlag(block)

	log.Println(blockInt.Int64())

	rpc, err := cmd.Flags().GetString("rpc")
	if err != nil {
		return err
	}

	snapshotPath, err := cmd.Flags().GetString("snapshot")
	if err != nil {
		return err
	}

	if snapshotPath != "" {
		snapshotBLock := validateBlockFlag(GetSnapshotBlockNumber(snapshotPath))
		if snapshotBLock.String() != blockInt.String() {
			return errors.New(fmt.Sprintf("snapshot block %s is NOT equeal to provided  block %s", snapshotBLock, blockInt))
		}
		log.Println(fmt.Sprintf("snapshot block %s is equeal to provided  block %s", snapshotBLock, blockInt))
		snapshotAddresses := ReadSnapshotData(snapshotPath)
		touchedAddresses := ReturnTouchedAddresses(FetchTracesForBlock(rpc, blockInt))
		seenAddresses := make([]common.Address, 0)
		unseenAddresses := make([]common.Address, 0)
		for _, v := range touchedAddresses {
			if snapshotAddresses[v] == true {
				seenAddresses = append(seenAddresses, v)
			} else {
				unseenAddresses = append(unseenAddresses, v)
			}
		}
		writeAddressesToFile("unseen.txt", unseenAddresses)
		writeAddressesToFile("seen.txt", seenAddresses)
	} else {
		traces := FetchTracesForBlock(rpc, blockInt)
		touchedAddresses := ReturnTouchedAddresses(traces)
		writeAddressesToFile("seen.txt", touchedAddresses)
	}
	return nil
}
