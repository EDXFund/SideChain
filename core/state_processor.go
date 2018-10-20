// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"github.com/EDXFund/Validator/common"
	"github.com/EDXFund/Validator/consensus"
	"github.com/EDXFund/Validator/core/state"
	"github.com/EDXFund/Validator/core/types"
	"github.com/EDXFund/Validator/core/vm"
//	"github.com/EDXFund/Validator/crypto"
	"github.com/EDXFund/Validator/params"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.ContractResults,  uint64, error) {
	var (
		receipts types.ContractResults
		usedGas  = new(uint64)
		header   = block.Header()
		//allLogs  []*types.Log
		//gp       = new(GasPool).AddGas(block.GasLimit())
	)


	/* MUST TODO, add contract and token id managements routine
	for i, tx := range block.Transactions() {

		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			//return nil, nil, 0, err
			rejections = append(rejections, &types.RejectInfo{TxHash: tx.Hash(), Reason: 1})
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	//if master chain process blcokinfos and rejections

	///// MUST TODO  blockInfos has already been included in block
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)

	*/
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(),receipts)


	return receipts, *usedGas, nil
}


func verifySigner() {
	//this is unnessary, because tx with invalid signature would not be placed into pool
}

func callSC(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.ContractResult, uint64, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, 0, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, 0, err
	}
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := types.NewContractResult(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	//if msg.To() == nil {
	//	receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	//}
	// Set the receipt logs and create a bloom for filtering
	//receipt.Logs = statedb.GetLogs(tx.Hash())
	//receipt.Bloom = types.CreateBloom(types.ContractResults{receipt})

	return receipt, gas, err
}
//if txType is 'T' : ApplyTransaction validate tx's signature,
//if TxType is 'C' : Synchronize data from master node and process data, then create contractResult
//if txType is "D" : it should be a token management tx
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.ContractResult, uint64, error) {

	switch tx.Type() {

	case types.TT_TRANSACTION:
		return nil,0,nil
	case types.TT_CONTRACT:
		return callSC(config , bc , author , gp , statedb , header , tx , usedGas, cfg )
	case types.TT_TOKEN:
		return nil,0,nil
	default:
		return nil,0,ErrNoGenesis
	}

}
