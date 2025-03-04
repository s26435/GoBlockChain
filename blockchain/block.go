package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	u "go.mod/utils"
)

type Block struct {
	Timestamp     int64
	Transactions []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}


func (b *Block)Serialize()[]byte{
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	u.Must(encoder.Encode(b))
	return result.Bytes()
}


func DeserializeBlock(d []byte)*Block{
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	u.Must(decoder.Decode(&block))
	return &block
}


func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(),transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock(coinbase *Transaction)*Block{
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block)HashTransactions()[]byte{
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions{
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

