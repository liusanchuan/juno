package core

import (
	"encoding/hex"

	"github.com/NethermindEth/juno/clients"
	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/bits-and-blooms/bitset"
)

type Block struct {
	// The hash of this block’s parent
	ParentHash *felt.Felt
	// The number (height) of this block
	Number uint64
	// The state commitment after this block
	GlobalStateRoot *felt.Felt
	// The StarkNet address of the sequencer who created this block
	SequencerAddress *felt.Felt
	// The time the sequencer created this block before executing transactions
	Timestamp uint64
	// The number of transactions in a block
	TransactionCount uint64
	// A commitment to the transactions included in the block
	TransactionCommitment *felt.Felt
	// The number of events
	EventCount uint64
	// A commitment to the events produced in this block
	EventCommitment *felt.Felt
	// The version of the StarkNet protocol used when creating this block
	ProtocolVersion uint64
	// Extraneous data that might be useful for running transactions
	ExtraData *felt.Felt
}

func (b *Block) Hash() *felt.Felt {
	// Todo: implement pedersen hash as defined here
	// https://docs.starknet.io/documentation/develop/Blocks/header/#block_hash
	return nil
}

type (
	commitmentStorage     map[string]string
	commitmentTrieStorage struct {
		storage commitmentStorage
	}
)

func (s *commitmentTrieStorage) Put(key *bitset.BitSet, value *TrieNode) error {
	keyEnc, err := key.MarshalBinary()
	if err != nil {
		return err
	}
	vEnc, err := value.MarshalBinary()
	if err != nil {
		return err
	}
	s.storage[hex.EncodeToString(keyEnc)] = hex.EncodeToString(vEnc)
	return nil
}

func (s *commitmentTrieStorage) Get(key *bitset.BitSet) (*TrieNode, error) {
	keyEnc, _ := key.MarshalBinary()
	value, found := s.storage[hex.EncodeToString(keyEnc)]
	if !found {
		panic("not found")
	}

	v := new(TrieNode)
	decoded, _ := hex.DecodeString(value)
	err := v.UnmarshalBinary(decoded)
	return v, err
}

func (s *commitmentTrieStorage) Delete(key *bitset.BitSet) error {
	keyEnc, _ := key.MarshalBinary()
	delete(s.storage, hex.EncodeToString(keyEnc))
	return nil
}

// TransactionCommitment is the root of a height 64 binary Merkle Patricia tree of the
// transaction hashes and signatures in a block.
func TransactionCommitment(transactions []*clients.Transaction) (*felt.Felt, error) {
	transactionCommitmentStorage := &commitmentTrieStorage{
		storage: make(commitmentStorage),
	}
	transactionCommitmentTrie := NewTrie(transactionCommitmentStorage, 64)

	zeroFelt := new(felt.Felt)
	emptySignatureHash, err := crypto.Pedersen(zeroFelt, zeroFelt)
	if err != nil {
		return nil, err
	}
	for i, transaction := range transactions {
		var signaturesHash *felt.Felt
		if transaction.Type == "INVOKE_FUNCTION" {
			signaturesHash, err = crypto.PedersenArray(transaction.Signature...)
			if err != nil {
				return nil, err
			}
		} else {
			signaturesHash = emptySignatureHash
		}
		transactionAndSignatureHash, err := crypto.Pedersen(transaction.Hash, signaturesHash)
		if err != nil {
			return nil, err
		}
		err = transactionCommitmentTrie.Put(new(felt.Felt).SetInt64(int64(i)), transactionAndSignatureHash)
		if err != nil {
			return nil, err
		}
	}

	return transactionCommitmentTrie.Root()
}

// EventCommitment is the root of a height 64 binary Merkle Patricia tree of the
// events in a block.
func EventCommitment(receipts []*clients.TransactionReceipt) (*felt.Felt, error) {
	eventCommitmentStorage := &commitmentTrieStorage{
		storage: make(commitmentStorage),
	}
	eventCommitmentTrie := NewTrie(eventCommitmentStorage, 64)

	var index int64
	for _, receipt := range receipts {
		for _, event := range receipt.Events {
			keys, err := crypto.PedersenArray(event.Keys...)
			if err != nil {
				return nil, err
			}

			data, err := crypto.PedersenArray(event.Data...)
			if err != nil {
				return nil, err
			}

			eventHash, err := crypto.PedersenArray(
				event.From,
				keys,
				data,
			)
			if err != nil {
				return nil, err
			}

			err = eventCommitmentTrie.Put(new(felt.Felt).SetInt64(index), eventHash)
			if err != nil {
				return nil, err
			}
			index++
		}
	}

	return eventCommitmentTrie.Root()
}
