package blockchain

import (
	"encoding/hex"
	"os"
	"fmt"

	"github.com/boltdb/bolt"
	u "go.mod/utils"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "Student - debil. Piwo to moje paliwo."

type BlockChain struct{
	tip []byte
	DB *bolt.DB
}
func (bc *BlockChain) Mine(transactions []*Transaction) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	u.Must(err)

	newBlock := NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		u.Must(err)

		err = b.Put([]byte("l"), newBlock.Hash)
		u.Must(err)
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewBlockChain(address string)*BlockChain{
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	u.Must(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil{
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			u.Must(err)
			err = b.Put(genesis.Hash, genesis.Serialize())
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

func (bc *BlockChain) FindUnspendTransactions(address string)[]Transaction{
	var unspentTX []Transaction

	spentTX := make(map[string][]int)
	bci:=bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions{
			txId := hex.EncodeToString(tx.ID)

			for outIdx, out  := range tx.Vout{
				if spentTX[txId] != nil{
					for _, spentOut := range  spentTX[txId]{
						if spentOut == outIdx{
							continue
						}
					}
				}
				if out.CanBeUnlockedWith(address){
					unspentTX = append(unspentTX, *tx)
				}
			}
			if !tx.IsCoinbase(){
				for _, in := range tx.Vin{
					if in.CanUnlockOutputWith(address){
						inTxID := hex.EncodeToString(in.Txid)
						spentTX[inTxID] = append(spentTX[inTxID], in.Vout)
					}
				}
			}
		}
		if len(block.PrevBlockHash) == 0{
			break
		}
	}
	return unspentTX
}


func (bc *BlockChain) FindUTXO(address string)[]TXOutput{
	var UTXO []TXOutput

	unspentTransactions := bc.FindUnspendTransactions(address)

	for _, tx := range unspentTransactions{
		for _, out := range tx.Vout{
			if out.CanBeUnlockedWith(address){
				UTXO = append(UTXO, out)
			}
		}
	}
	return UTXO
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}


func CreateBlockchain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	u.Must(err)

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		u.Must(err)
		err = b.Put(genesis.Hash, genesis.Serialize())
		u.Must(err)

		err = b.Put([]byte("l"), genesis.Hash)
		u.Must(err)
		tip = genesis.Hash

		return nil
	})

	u.Must(err)

	bc := BlockChain{tip, db}

	return &bc
}