package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	//Data          []byte
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func Deserialize(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		log.Panic(err)
	}

	return &block
}

//func (b *Block) SetHash() {
//	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//	key := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
//	hash := sha256.Sum256(key)
//
//	b.Hash = hash[:]
//}

//func NewBlock(data string, prevBlockHash []byte) *Block {
//	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
//
//	pow := NewProofOfWork(block)
//	nonce, hash := pow.Run()
//
//	block.Hash = hash
//	block.Nonce = nonce
//
//	return block
//}

func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), txs, prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (block *Block) HashTransactions() []byte {
	var txHashs [][]byte

	for _, tx := range block.Transactions {
		txHashs = append(txHashs, tx.ID)
	}

	hash := sha256.Sum256(bytes.Join(txHashs, []byte{}))

	return hash[:]
}
