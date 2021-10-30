package main

import (
	"github.com/gin-gonic/gin"
	"github.com/streamdp/goblockchain/internal/bchain"
	"log"
	"net/http"
	"strconv"
)

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
