package blockchain

import (
	"github.com/boltdb/bolt"

	u "go.mod/utils"
)
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainIterator)Next()*Block{
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	u.Must(err)

	i.currentHash = block.PrevBlockHash
	return block
}