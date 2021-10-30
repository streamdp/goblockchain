package bchain

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//InitBlockChain - make new blockchain
func InitBlockChain() (bChain *BlockChain) {
	bChain = &BlockChain{}
	bChain.Complexity = 5
	bChain.NewBlock(1)
	return bChain
}

//PingNode - check if the node is alive
func PingNode(url string) bool {
	var body []byte
	var err error
	type jsonPing struct {
		Message string `json:"message"`
	}
	if body, err = GetBodyResponse(url + "ping"); err != nil {
		log.Println(err)
	}
	answer := &jsonPing{}
	if err = json.Unmarshal(body, answer); err != nil {
		log.Println(err)
		return false
	}
	if answer.Message == "pong" {
		return true
	}
	return false
}

//GetBodyResponse - get response body by url line
func GetBodyResponse(url string) ([]byte, error) {
	var body []byte
	var err error
	var req *http.Request
	var resp *http.Response
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	return body, nil
}

//GetChain - get blockchain by url line from another node
func GetChain(url string) ([]*Block, error) {
	var body []byte
	var err error
	type jsonChain struct {
		Chain  []*Block `json:"chain"`
		Length int64    `json:"-"`
	}
	if body, err = GetBodyResponse(url + "chain"); err != nil {
		return nil, err
	}
	chain := &jsonChain{}
	if err = json.Unmarshal(body, chain); err != nil {
		return nil, err
	}
	return chain.Chain, nil
}
