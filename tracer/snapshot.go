package tracer

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"

	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func ReadSnapshotData(path string) map[common.Address]bool {
	log.Println("Start parsing snapshot")
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	addresses := make(map[common.Address]bool)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	maxBufferSize := 1024 * 1024 * 2048
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	for scanner.Scan() {
		acc := new(state.DumpAccount)
		if err := json.Unmarshal(scanner.Bytes(), acc); err != nil {
			continue
		}
		if acc.Address != nil {
			addresses[*acc.Address] = true
			if acc.Address.String() == "0x472370dd80620cd6d3b8FD57daDc0B9752988e27" {
				fmt.Println("0x472370dd80620cd6d3b8FD57daDc0B9752988e27 found in snapshot")
			}
		}
	}
	log.Printf("Read snapshot finished with %s addresses", len(addresses))
	return addresses
}

func GetSnapshotBlockNumber(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	maxBufferSize := 64 * 2048
	buf := make([]byte, 0, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)
	for i := 0; i < 30; i += 1 {
		scanner.Scan()
		if strings.Contains(scanner.Text(), "State dump configured") {
			arr := strings.Split(scanner.Text(), " ")
			for _, v := range arr {
				if strings.Contains(v, "block") {
					arr = strings.Split(v, "=")
					block := strings.Replace(arr[1], ",", "", -1)
					return block
				}
			}
			break
		}
	}
	log.Fatal("Snapshot does not contain block number")
	return ""
}
