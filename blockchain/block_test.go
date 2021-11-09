package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGenerateNewBlock(t *testing.T) {
	genesisBlock := GenerateGenesisBlock()

	amount := rand.Intn(1000)

	// Actual block
	actualBlock, err := GenerateBlock(genesisBlock, amount)
	if err != nil {
		t.Error(err)
	}

	// Expected block
	expectedBlock := Block{
		Index:     genesisBlock.Index + 1,
		Timestamp: actualBlock.Timestamp, // Assign the same time as the generated block
		Amount:    amount,
		Hash:      "",
		PrevHash:  genesisBlock.Hash,
	}
	expectedBlock.Hash = CalculateHash(expectedBlock)

	if !reflect.DeepEqual(actualBlock, expectedBlock) {
		t.Error("Expected s, got s")
	}
}

func TestCalculateHash(t *testing.T) {
	genesisBlock := GenerateGenesisBlock()

	block := Block{
		Index:     rand.Intn(1000),
		Timestamp: time.Now().String(),
		Amount:    rand.Intn(1000),
		Hash:      "",
		PrevHash:  genesisBlock.Hash,
	}

	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Amount) + block.PrevHash
	sha256Gen := sha256.New()
	sha256Gen.Write([]byte(record))
	actualHash := hex.EncodeToString(sha256Gen.Sum(nil))

	expectedHash := CalculateHash(block)

	if expectedHash != actualHash {
		t.Errorf("Expected %s, got %s", expectedHash, actualHash)
	}
}
