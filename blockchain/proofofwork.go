package blockchain

import(
	"math/big"
	u "go.mod/utils"
	"math"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const targetBits = 24
var (
	maxNonce = math.MaxInt64
)


type ProofOfWork struct{
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block)*ProofOfWork{
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork)prepareData(nonce int)[]byte{
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			u.IntToHex(pow.block.Timestamp),
			u.IntToHex(int64(targetBits)),
			u.IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork)Run()(int, []byte){
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining new Block")
	for nonce < maxNonce{
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1{
			break
		}else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork)Validate() bool{
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
