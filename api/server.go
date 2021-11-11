package api

import (
	"fmt"
	"net/http"

	"github.com/andregri/tiny-p2p-blockchain/blockchain"
	"github.com/gin-gonic/gin"
)

func Start(port int) {
	router := gin.Default()
	router.GET("/blockchain", getBlockchain)

	router.Run(fmt.Sprintf("localhost:%d", port))
}

func getBlockchain(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, blockchain.Blockchain)
}
