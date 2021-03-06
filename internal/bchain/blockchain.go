package bchain

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// BlockChain - basic blockchain structure
type BlockChain struct {
	Chain              []*Block       `json:"chain"`
	CurrentTransaction []*Transaction `json:"current_transaction"`
	Nodes              []*url.URL     `json:"nodes"`
	Complexity         int            `json:"complexity"`
	TimeStamp          int64          `json:"time_stamp"`
}

// Block - basic block structure
type Block struct {
	Index        int64          `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	Proof        int64          `json:"proof"`
	PreviousHash string         `json:"previous_hash"`
}

// Transaction - basic transaction structure
type Transaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int64  `json:"amount"`
}

// UpdateTimeStamp - update TimeStamp for any change in the blockchain
func (b *BlockChain) UpdateTimeStamp() {
	b.TimeStamp = time.Now().Unix()
}

// IncreaseComplexity - increases complexity of mining a block if necessary
func (b *BlockChain) IncreaseComplexity() {
	if b.Complexity < 38 {
		b.Complexity++
		b.UpdateTimeStamp()
	}
}

// DecreaseComplexity - decreases complexity of mining a block if necessary
func (b *BlockChain) DecreaseComplexity() {
	if b.Complexity > 4 {
		b.Complexity--
		b.UpdateTimeStamp()
	}
}

// RegisterNode - add new node to node list
func (b *BlockChain) RegisterNode(urlString string) bool {
	nodeUrl, err := url.Parse(urlString)
	if err != nil {
		log.Println(err)
	}
	for _, knownUrl := range b.Nodes {
		if *knownUrl == *nodeUrl {
			return false
		}
	}
	if PingNode(nodeUrl.String()) {
		b.Nodes = append(b.Nodes, nodeUrl)
		b.UpdateTimeStamp()
		return true
	}
	return false
}

// NewBlock - make new block and add to blockchain
func (b *BlockChain) NewBlock(proof int64) (block *Block) {
	var index int64
	var previousHash string
	if len(b.Chain) > 0 {
		index = b.LastBlock().Index + 1
		previousHash = b.Hash(b.LastBlock())
	} else {
		previousHash = "100"
		index = int64(0)
	}
	b.Chain = append(b.Chain, &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: b.CurrentTransaction,
		Proof:        proof,
		PreviousHash: previousHash,
	})
	b.CurrentTransaction = []*Transaction{}
	b.UpdateTimeStamp()
	return b.LastBlock()
}

// NewTransaction - make new transaction and add to blockchain
func (b *BlockChain) NewTransaction(sender, recipient string, amount int64) (idx int64) {
	b.CurrentTransaction = append(b.CurrentTransaction, &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	})
	b.UpdateTimeStamp()
	return b.LastBlock().Index + 1
}

// Hash - calculate SHA1 hash of the block
func (b *BlockChain) Hash(block *Block) (hashString string) {
	blockString, err := json.Marshal(&block)
	if err != nil {
		log.Println(err)
	}
	hash := sha1.New()
	hash.Write(blockString)
	return hex.EncodeToString(hash.Sum(nil))
}

// LastBlock - return last block of the blockchain
func (b *BlockChain) LastBlock() (block *Block) {
	return b.Chain[len(b.Chain)-1]
}

// ProofOfWork - do some block mining work
func (b *BlockChain) ProofOfWork(lastProof int64) (proof int64) {
	for !b.ValidProof(lastProof, proof) {
		proof++
	}
	return proof
}

// ValidProof - see if the proof is right
func (b *BlockChain) ValidProof(lastProof, proof int64) (valid bool) {
	guess := strconv.FormatInt(lastProof, 10) + strconv.FormatInt(proof, 10)
	hash := sha1.New()
	hash.Write([]byte(guess))
	guessHash := hex.EncodeToString(hash.Sum(nil))
	return guessHash[:b.Complexity] == strings.Repeat("0", b.Complexity)
}

// IsValidChain - integrity check of the blockchain
func (b *BlockChain) IsValidChain(chain []*Block) bool {
	previousBlock := chain[0]
	for i := 1; i < len(chain); i++ {
		currentBlock := chain[i]
		if currentBlock.PreviousHash != b.Hash(previousBlock) {
			return false
		}
		if !b.ValidProof(previousBlock.Proof, currentBlock.Proof) {
			return false
		}
		previousBlock = currentBlock
	}
	return true
}

// ResolveConflicts - replace blockchain if if conflict
func (b *BlockChain) ResolveConflicts() bool {
	var chain []*Block
	var err error
	var actualChain []*Block
	maxLength := len(b.Chain)
	for _, node := range b.Nodes {
		if chain, err = GetChain(node.String()); err != nil {
			log.Println("error with getting chain from node", err)
		}
		if len(chain) > maxLength && b.IsValidChain(chain) {
			maxLength = len(chain)
			actualChain = chain
		}
	}
	if len(actualChain) > 0 {
		b.Chain = actualChain
		b.UpdateTimeStamp()
		return true
	}
	return false
}
