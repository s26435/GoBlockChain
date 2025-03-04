package blockchain

import (
	"github.com/boltdb/bolt"
	u "go.mod/utils"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type BlockChain struct{
	tip []byte
	DB *bolt.DB
}

func (bc *BlockChain) AddBlock(data string){
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	u.Must(err)

	newBlock := NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		u.Must(err)

		err = b.Put([]byte("l"), newBlock.Hash)
		u.Must(err)

		bc.tip = newBlock.Hash
		return nil
	})
	u.Must(err)
}

func NewBlockChain()*BlockChain{
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	u.Must(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil{
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			u.Must(err)
			u.Must(b.Put(genesis.Hash, genesis.Serialize()))
			u.Must(b.Put([]byte("l"), genesis.Hash))
			tip = genesis.Hash
		}else{
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	u.Must(err)
	bc := BlockChain{tip, db}
	return &bc
}


func (bc *BlockChain)Iterator()*BlockchainIterator{
	bci := &BlockchainIterator{bc.tip, bc.DB}
	return bci
}