package core

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/NethermindEth/juno/clients"
)

var (
	//go:embed testdata/block_156000.json
	transactions156000 []byte

	//go:embed testdata/block_1.json
	transactions1 []byte

	//go:embed testdata/block_1_integration.json
	transactions1Integration []byte

	//go:embed testdata/block_16789_main.json
	transactions16789Main []byte
)

func getBlocks(t *testing.T) []*clients.Block {
	// https://alpha4.starknet.io/feeder_gateway/get_block?blockNumber=156000
	var txs156000 *clients.Block
	err := json.Unmarshal(transactions156000, &txs156000)
	if err != nil {
		t.Fatal(err)
	}
	// https://alpha4.starknet.io/feeder_gateway/get_block?blockNumber=1
	var txs1 *clients.Block
	err = json.Unmarshal(transactions1, &txs1)
	if err != nil {
		t.Fatal(err)
	}
	// https://external.integration.starknet.io/feeder_gateway/get_block?blockNumber=1
	var txs1Integration *clients.Block
	err = json.Unmarshal(transactions1Integration, &txs1Integration)
	if err != nil {
		t.Fatal(err)
	}
	// https://alpha-mainnet.starknet.io/feeder_gateway/get_block?blockNumber=16789
	var txs16789Main *clients.Block
	err = json.Unmarshal(transactions16789Main, &txs16789Main)
	if err != nil {
		t.Fatal(err)
	}

	return []*clients.Block{txs156000, txs1, txs1Integration, txs16789Main}
}

func TestTransactionCommitment(t *testing.T) {
	blocks := getBlocks(t)
	tests := []struct {
		txs  []*clients.Transaction
		want string
	}{
		{
			blocks[0].Transactions,
			"0x24638e0ca122d0260d54e901dc0942ea68bd1fc40a96b5da765985c47c92500",
		},
		{
			blocks[1].Transactions,
			"0x18bb7d6c1c558aa0a025f08a7d723a44b13008ffb444c432077f319a7f4897c",
		},
		{
			blocks[2].Transactions,
			"0xbf11745df434cbd284e13ca36354139a4bca2f6722e737c6136590990c8619",
		},
		// TODO: Fix this failing test
		// {
		// 	blocks[3].Transactions,
		// 	"0x580a06bfc8c3fe39bbb7c5d16298b8928bf7c28f4c31b8e6b48fc25cd644fc1",
		// },
	}

	for _, test := range tests {
		commitment, _ := ComputeTransactionCommitment(test.txs)
		if "0x"+commitment.Text(16) != test.want {
			t.Errorf("got %s, want %s", commitment, test.want)
		}
	}
}

func TestEventCommitment(t *testing.T) {
	blocks := getBlocks(t)
	tests := []struct {
		receipts []*clients.TransactionReceipt
		want     string
	}{
		{
			blocks[0].Receipts,
			"0x5d25e41d43b00681cc63ed4e13a82efe3e02f47e03173efbd737dd52ba88c7e",
		},
		{
			blocks[1].Receipts,
			"0x0",
		},
		{
			blocks[2].Receipts,
			"0x0",
		},
		{
			blocks[3].Receipts,
			"0x6f499789aabb31935810ce89d6ea9e9d37c5921c0d7fae2bd68f2fff5b7b93f",
		},
	}

	for _, test := range tests {
		commitment, _ := ComputeEventCommitment(test.receipts)
		if "0x"+commitment.Text(16) != test.want {
			t.Errorf("got %s, want %s", "0x"+commitment.Text(16), test.want)
		}
	}
}
