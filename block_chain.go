package main

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "block_chain.db"
const blocksBucket = "blocks"
const LastHashKey = "l"

type Blockchain struct {
	id []byte
	db *bolt.DB
}

//func (bc *Blockchain) AddBlock(data string) {
//	//prevBlock := bc.blocks[len(bc.blocks)-1]
//	//block := NewBlock(data, prevBlock.Hash)
//	//bc.blocks = append(bc.blocks, block)
//	var lastHash []byte
//	err := bc.db.View(func(tx *bolt.Tx) error {
//		b := tx.Bucket([]byte(blocksBucket))
//		lastHash = b.Get([]byte(LastHashKey))
//
//		return nil
//	})
//
//	if err != nil {
//		log.Panic(err)
//	}
//
//	block := NewBlock(data, lastHash)
//
//	bc.addBlock(block)
//	bc.id = block.Hash
//}

func (bc *Blockchain) AddBlock(block *Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(block.Hash, block.Serialize())
		err = b.Put([]byte(LastHashKey), block.Hash)

		if err != nil {
			log.Panic(err)
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (bc *Blockchain) MineBlock(txs []*Transaction) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		lashHash := b.Get([]byte(LastHashKey))

		newBlock := NewBlock(txs, lashHash)

		b.Put([]byte(newBlock.Hash), newBlock.Serialize())
		b.Put([]byte(LastHashKey), newBlock.Hash)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func NewBlockchain(address string) *Blockchain {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var hash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			cbTx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbTx)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte(LastHashKey), genesis.Hash)

			if err != nil {
				log.Panic(err)
			}

			hash = genesis.Hash

		} else {
			hash = b.Get([]byte(LastHashKey))
		}

		return nil

	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{hash, db}

}

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {

	unspentOutputs := make(map[string][]int)

	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {

			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

			}

			if accumulated >= amount {
				break Work
			}

		}

	}

	return accumulated, unspentOutputs
}
