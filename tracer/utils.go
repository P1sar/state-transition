package tracer

import (
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"os"
	"strings"
)

func validateBlockFlag(block string) *big.Int {
	v := big.NewInt(0)
	if strings.HasPrefix(block, "0x") || strings.HasPrefix(block, "0X") {
		log.Println("HERE")
		v.SetString(block, 0)
		return v
	}
	v.SetString(block, 10)
	return v
}

func writeAddressesToFile(name string, addresses []common.Address) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, v := range addresses {
		f.WriteString(v.String())
		f.WriteString("\n")
	}

}
