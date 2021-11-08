module github.com/andregri/tiny-p2p-blockchain

go 1.16

require (
	github.com/libp2p/go-libp2p v0.15.1
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.11.0
	github.com/libp2p/go-libp2p-kad-dht v0.15.0
	github.com/libp2p/go-libp2p-noise v0.3.0
	github.com/libp2p/go-libp2p-tls v0.3.0
)

replace github.com/libp2p/go-libp2p => ./go-libp2p
