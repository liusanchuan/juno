package main

import (
	"math"
	"strconv"

	"github.com/wealdtech/go-merkletree/keccak256"
)

// Contract Definition Structure as defined here
// https://docs.starknet.io/documentation/develop/Contracts/contract-hash/
// https://www.cairo-lang.org/docs/hello_starknet/l1l2.html#receiving-a-message-from-l1
type ContractDefinition struct {
	// Need to change the field types to felt
	api_version          string
	external_functions   ContractEntryPoint
	l1_handlers          ContractEntryPoint
	constructors         ContractEntryPoint
	builtin_list         []byte
	bytecode             []byte
	entry_points_by_type map[string]int
}

type ContractEntryPoint struct {
	selector string `json:"selector"`
	offset   string `json:"offset"`
}

type Program struct {
	prime             int
	data              []int
	hints             map[int]string
	builtins          []string
	main_scope        string
	identifiers       string
	reference_manager string
	attributes        []int
	debug_info        bool
}

func main() {
	contract_definition := &ContractDefinition{}
	compute_class_hash(contract_definition)
}

// Compute the class hash for a given contract definition
// Class hash definition is mentioned here https://docs.starknet.io/documentation/develop/Contracts/contract-hash/
// Reference taken from the cairo implementation
// https://github.com/starkware-libs/cairo-lang/blob/7712b21fc3b1cb02321a58d0c0579f5370147a8b/src/starkware/starknet/core/os/contracts.cairo#L47
func compute_class_hash(contract_definition *ContractDefinition) []byte {
	// def, err := json.Marshal(contract_definition)
	// if err != nil {
	// 	fmt.Println("Error while marshaling the contract definition ", err)
	// }
	hash_state := pedersen_hash(0, contract_definition.api_version)
	hash_state = pedersen_hash(hash_state, contract_definition.external_functions)
	hash_state = pedersen_hash(hash_state, contract_definition.l1_handlers)
	hash_state = pedersen_hash(hash_state, contract_definition.constructors)
	hash_state = pedersen_hash(hash_state, contract_definition.builtin_list)
	hash_state = pedersen_hash(hash_state, contract_definition.bytecode)
	return hash_state
}

// Pedersen hash function
func pedersen_hash(input1 interface{}, input2 interface{}) []byte {
	return nil
}

// Function compute the starknet_keccak defined here
// https://docs.starknet.io/documentation/develop/Hashing/hash-functions/#starknet_keccak
func starknet_keccak(data []byte) int {
	MASK_250 := int(math.Pow(2, 250) - 1)
	keccak := keccak256.New()
	hashed_data := keccak.Hash(data)
	result, _ := strconv.Atoi(string(hashed_data))
	return result & MASK_250
}
