package blockchain

import (
	"fmt"
	u "go.mod/utils"
	"bytes"
	"encoding/gob"
	"crypto/sha256"

)

var subsidy int = 1 //nagroda za wykopanie - co 210 000 bloków wykopanych jest zmniejszana o połowę

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward %s", to)
	}
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	u.Must(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}


func (txi *TXInput)CanUnlockOutputWith(unlockingData string)bool{
	return txi.ScriptSig == unlockingData
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (txo *TXOutput)CanBeUnlockedWith(unlockingData string) bool{
	return txo.ScriptPubKey == unlockingData
}