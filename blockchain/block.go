package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"

	u "go.mod/utils"
)

type Block struct {
	Timestamp     int64
	Data          []byte
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


func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock()*Block{
	return NewBlock("Genesis Block", []byte{})
}