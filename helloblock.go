package helloblock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HelloBlockTransaction struct {
	TxHash   string `json:"txHash"`
	RawTxHex string `json:"rawTxHex"`
}

type HelloBlockPropagateData struct {
	Transaction HelloBlockTransaction `json:"transaction"`
}

type HelloBlockPropagateReply struct {
	Status  string                  `json:"status"`
	Message string                  `json:"message"`
	Details []string                `json:"details"`
	Data    HelloBlockPropagateData `json:"data"`
}

type HelloBlockNetwork string

func (n *HelloBlockNetwork) String() string {
	return string(*n)
}

const (
	Testnet HelloBlockNetwork = "testnet"
	Mainnet                   = "mainnet"
)

const (
	success = "success"
)

var network = Testnet

func SetNetwork(newNetwork HelloBlockNetwork) {
	network = newNetwork
}

func Propagate(rawTxHex string) (*HelloBlockPropagateData, error) {
	u := fmt.Sprintf("https://%v.helloblock.io/v1/transactions/", network.String())
	resp, err := http.PostForm(u, url.Values{"rawTxHex": {rawTxHex}})
	if err != nil {
		return nil, fmt.Errorf("error calling %v: %v", u, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var reply HelloBlockPropagateReply
	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("could not process JSON %v: %v", body, err)
	}
	if reply.Status != success {
		return nil, fmt.Errorf("error from helloblock.io: %v: %v", reply.Message, strings.Join(reply.Details, ", "))
	}
	return &reply.Data, nil
}
