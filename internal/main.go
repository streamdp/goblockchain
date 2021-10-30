package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"goblockchain/internal/bchain"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const BlockReward = 1

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

func getNodeIdentifier() string {
	nodeUuid, _ := uuid.NewUUID()
	nodeIdentifier := nodeUuid.String()
	return strings.Replace(nodeIdentifier, "-", "", -1)
}

func increase(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		bc.IncreaseComplexity()
		c.JSON(http.StatusOK, gin.H{
			"message":    "Complexity successfully increased",
			"complexity": bc.Complexity,
		})
	}
}

func decrease(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		bc.DecreaseComplexity()
		c.JSON(http.StatusOK, gin.H{
			"message":    "Complexity successfully decreased",
			"complexity": bc.Complexity,
		})
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func mine(bc *bchain.BlockChain, nodeId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		proof := bc.ProofOfWork(bc.LastBlock().Proof)
		bc.NewTransaction("0", nodeId, BlockReward)
		block := bc.NewBlock(proof)
		c.JSON(http.StatusOK, gin.H{
			"message": "New Block Forged",
			"block":   block,
		})
	}
}

func transactionsNew(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := bchain.Transaction{}
		if err := c.BindJSON(&query); err != nil {
			log.Println(err)
		}
		if validTransaction(&query) {
			index := bc.NewTransaction(query.Sender, query.Recipient, query.Amount)
			c.JSON(http.StatusCreated, gin.H{
				"message": "Transaction will be added to Block " + strconv.FormatInt(index, 10),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Missing values",
			})
		}
	}
}

func validTransaction(t *bchain.Transaction) bool {
	if t.Amount > 0 && t.Sender != "" && t.Recipient != "" {
		return true
	}
	return false
}

func chain(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"chain":  bc.Chain,
			"length": len(bc.Chain),
		})
	}
}

func registerNodes(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		var messages []string
		type jsonQuery struct {
			Urls []string `json:"nodes"`
		}
		query := &jsonQuery{}
		if err := c.BindJSON(&query); err != nil {
			log.Println(err)
		}
		if len(query.Urls) > 0 {
			for i, node := range query.Urls {
				if bc.RegisterNode(node) {
					messages = append(messages, "New nodes have been added, ", strconv.Itoa(i), node)
				}
			}
		}
		if len(messages) == 0 {
			messages = append(messages, "No nodes was added (address already exist or not available)")
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":     messages,
			"total_nodes": len(bc.Nodes),
		})
	}
}

func consensus(bc *bchain.BlockChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		if bc.ResolveConflicts() {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Our chain was replaced",
				"chain":   bc.Chain,
			})
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"message": "Our chain is authoritative",
				"chain":   bc.Chain,
			})
		}
	}
}
