package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	Index     int // Position of the data record in the blockchain
	Timestamp string
	Amount    int    // amount
	Hash      string // SHA256 identifier of this data record
	PrevHash  string // sha256 identifier of the previous data record
}

var Blockchain []Block

func CalculateHash(block Block) string {
	// Concatenate index, timestamp bpm and prevhash to generate the block hash
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Amount) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString((hashed))
}

func GenerateBlock(oldBlock Block, amount int) (Block, error) {
	// set block variables
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Amount = amount
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}

func GenerateGenesisBlock() Block {
	t := time.Now()
	block := Block{0, t.String(), 0, "", " "}
	return block
}

func IsBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}
