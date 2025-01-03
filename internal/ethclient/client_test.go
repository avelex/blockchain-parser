package ethclient_test

import (
	"context"
	"testing"

	"github.com/avelex/blockchain-parser/internal/ethclient"
)

var testBlockchainRPC = "https://eth.llamarpc.com"

func Test_BlockNumber(t *testing.T) {
	client := ethclient.New(testBlockchainRPC)

	testCases := []struct {
		desc string
	}{
		{
			desc: "Latest Block Number",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			bn, err := client.BlockNumber(context.Background())
			if err != nil {
				t.Fatalf("failed to get block number: %v", err)
			}

			if bn == 0 {
				t.Fatalf("block number is empty")
			}
		})
	}
}

func Test_HeaderByNumber(t *testing.T) {
	client := ethclient.New(testBlockchainRPC)

	testCases := []struct {
		desc   string
		number int
		wantTx int
	}{
		{
			desc:   "Existing Block",
			number: 21543920,
			wantTx: 215,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			block, err := client.BlockHeaderByNumber(context.Background(), tC.number)
			if err != nil {
				t.Fatalf("failed to get block number: %v", err)
			}

			if block == nil {
				t.Fatalf("block is empty")
			}

			if block.Number != tC.number {
				t.Fatalf("block number is not equal to %d", tC.number)
			}
		})
	}
}

func Test_TransactionReceipt(t *testing.T) {
	client := ethclient.New(testBlockchainRPC)

	testCases := []struct {
		desc             string
		txHash           string
		wantFrom, wantTo string
	}{
		{
			desc:     "Existing Transaction",
			txHash:   "0xa095ab2eadeb8451e5eadc2329c8dbcabfae81bfd0b05d2f7c7fa635889b959b",
			wantFrom: "0xfe556e4f848c82093d0a33cc41761d18f67099ca",
			wantTo:   "0x22a7a914cf352f7361c199188a23da94fe71b277",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			receipt, err := client.TransactionReceipt(context.Background(), tC.txHash)
			if err != nil {
				t.Fatalf("failed to get transaction receipt: %v", err)
			}

			if receipt == nil {
				t.Fatalf("transaction receipt is empty")
			}

			if receipt.Hash != tC.txHash {
				t.Fatalf("transaction hash is not equal, want %s, got %s", tC.txHash, receipt.Hash)
			}

			if receipt.From != tC.wantFrom {
				t.Fatalf("transaction from is not equal, want %s, got %s", tC.wantFrom, receipt.From)
			}

			if receipt.To != tC.wantTo {
				t.Fatalf("transaction to is not equal, want %s, got %s", tC.wantTo, receipt.To)
			}
		})
	}
}
