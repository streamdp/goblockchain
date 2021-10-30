package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streamdp/goblockchain/internal/bchain"
	"strconv"
	"strings"
)

// BlockReward - block mining reward
const BlockReward = 1

func getNodeIdentifier() string {
	nodeUuid, _ := uuid.NewUUID()
	nodeIdentifier := nodeUuid.String()
	return strings.Replace(nodeIdentifier, "-", "", -1)
}

func main() {
	var showHelp bool
	var port int

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.IntVar(&port, "port", 8080, "set specify port")
	flag.Parse()
	if showHelp {
		fmt.Println("goBlockChain is a version of the classic blockChain with POW algorithms \n" +
			"written in golang. May come in handy when learning about this technology.")
		fmt.Println("")
		flag.Usage()
		return
	}
	nodeIdentifier := getNodeIdentifier()
	bChain := bchain.InitBlockChain()

	r := gin.Default()
	r.GET("/ping", ping)
	r.GET("/mine", mine(bChain, nodeIdentifier))
	r.GET("/chain", chain(bChain))
	r.GET("/nodes/resolve", consensus(bChain))
	r.GET("/mine/complexity/increase", increase(bChain))
	r.GET("/mine/complexity/decrease", decrease(bChain))
	r.POST("/transactions/new", transactionsNew(bChain))
	r.POST("/nodes/register", registerNodes(bChain))

	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		return
	}
}
