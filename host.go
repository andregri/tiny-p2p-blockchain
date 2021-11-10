package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"

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

	peerChan := peer.InitMdns(host, "rendesvouz")

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
