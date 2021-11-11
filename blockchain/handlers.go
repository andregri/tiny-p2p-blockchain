package blockchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Synchronize functions that add new blocks
var mutex = &sync.Mutex{}

func WriteBlockchain(rw *bufio.ReadWriter, blockChannel <-chan Block) {
	go broadcastChain(rw)

	for {
		// Wait a new block
		newBlock := <-blockChannel

		mutex.Lock()

		// Add to chain
		if IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			Blockchain = append(Blockchain, newBlock)
		}

		// Convert to json
		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}

		// Send blockchain
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()

		mutex.Unlock()
	}
}

func ReadBlockchain(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Receive chain
			chain := make([]Block, 0)
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()

			// If received chain is longer, drop old chain
			if len(chain) > len(Blockchain) {
				Blockchain = chain
			}
			bytes, err := json.MarshalIndent(Blockchain, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			const green string = "\x1b[32m"
			const resetColor string = "\x1b[0m"
			fmt.Printf("%s%s%s>", green, string(bytes), resetColor)

			mutex.Unlock()
		}
	}
}

func broadcastChain(rw *bufio.ReadWriter) {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}
}
