module github.com/andregri/tiny-p2p-blockchain

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/libp2p/go-libp2p v0.15.1
	github.com/libp2p/go-libp2p-core v0.11.0
	github.com/multiformats/go-multiaddr v0.4.1
)

replace github.com/libp2p/go-libp2p => ./go-libp2p
