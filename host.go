package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andregri/tiny-p2p-blockchain/blockchain"
	"github.com/andregri/tiny-p2p-blockchain/peer"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
)

var blockChannel = make(chan blockchain.Block)

func main() {
	listenAddr := flag.String("listen", "", "listen address")
	listenPort := flag.Int("port", 0, "listen port")
	proto := flag.String("proto", "blk", "protocol name")
	flag.Parse()

	if *listenPort == 0 {
		panic("port must be > 0")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	ctx := context.Background()

	host, _ := peer.GenerateNewHost(*listenAddr, *listenPort)

	// Set a function as stream handler.
	// This function is called when a peer initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(protocol.ID(*proto), handleStream)

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", *listenAddr, *listenPort, host.ID().Pretty())

	// Discover new peers
	peerChan := peer.InitMdns(host, "rendesvouz")

	// Add genesis block if chain is empty
	go func() {
		if len(blockchain.Blockchain) < 1 {
			genesis := blockchain.GenerateGenesisBlock()
			blockchain.Blockchain = append(blockchain.Blockchain, genesis)
		}

		generateRandomBlocks()
	}()

	// Connect to new peers
	for {
		peer := <-peerChan // will block untill we discover a peer
		fmt.Println("Found peer:", peer, ", connecting")

		if err := host.Connect(ctx, peer); err != nil {
			log.Println("Connection failed:", err)
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(*proto))

		if err != nil {
			log.Println("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go blockchain.WriteBlockchain(rw, blockChannel)
			go blockchain.ReadBlockchain(rw)
			fmt.Println("Connected to:", peer)
		}
	}
}

func handleStream(stream network.Stream) {
	fmt.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go blockchain.WriteBlockchain(rw, blockChannel)
	go blockchain.ReadBlockchain(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func generateRandomBlocks() {
	for {
		// Sleep for a random time
		sleepTime := rand.Intn(20) + 10
		time.Sleep(time.Duration(sleepTime) * time.Second)

		// Generate a new block
		amount := rand.Intn(1000)
		lastBlockIndex := len(blockchain.Blockchain) - 1
		oldBlock := blockchain.Blockchain[lastBlockIndex]
		block, err := blockchain.GenerateBlock(oldBlock, amount)
		if err != nil {
			log.Println(err)
		}

		log.Printf("Block %d generated\n", amount)
		blockChannel <- block
	}
}
