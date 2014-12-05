package helloblock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HelloBlockAddress struct {
	Balance uint64 `json:"balance"`
	ConfirmedBalance uint64 `json:"confirmedBalance"`
	TxsCount uint `json:"txsCount"`
	ConfirmedTxsCount uint `json:"confirmedTxsCount"`
	TotalReceivedValue uint64 `json:"totalReceivedValue"`
	ConfirmedReceivedValue uint64 `json:"confirmedReceivedValue"`
	Address string `json:"address"`
	Hash160 string `json:"hash160"`
	Type string `json:"type"`
}

type HelloBlockGetAddressData struct {
	Address HelloBlockAddress `json:"address"`
}

type HelloBlockReply struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Details []string             `json:"details"`
	DataRaw    json.RawMessage `json:"data"`
}

type HelloBlockUnspent struct {
	TxHash       string `json:"txHash"`
	Index        int    `json:"index"`
	ScriptPubKey string `json:"scriptPubKey"`
	Value        int    `json:"value"`
	Address      string `json:"address"`
}

type HelloBlockFaucetData struct {
	PrivateKeyWIF string              `json:"privateKeyWIF"`
	PrivateKeyHex string              `json:"privateKeyHex"`
	Address       string              `json:"address"`
	Hash160       string              `json:"hash160"`
	FaucetType    int                 `json:"faucetType"`
	Unspents      []HelloBlockUnspent `json:"unspents"`
}

type HelloBlockFaucetReply struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Details []string             `json:"details"`
	Data    HelloBlockFaucetData `json:"data"`
}

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

func GetAddress(address string) (*HelloBlockGetAddressData, error) {
	u := fmt.Sprintf("%s/%s", helloBlockURL("v1/addresses"), address)
	reply, err := getReply(u)
	if err != nil {
		return nil, err
	}
	var data HelloBlockGetAddressData
	if err := json.Unmarshal(reply.DataRaw, &data); err != nil {
		return nil, fmt.Errorf("couldn't process JSON %v: %v", reply.DataRaw, err)
	}
	return &data, nil
}

func getReply(u string) (*HelloBlockReply, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error calling %v: %v", u, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var reply HelloBlockReply
	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("could not process JSON %v: %v", string(body), err)
	}
	if reply.Status != success {
		return nil, fmt.Errorf("error from %v: %v: %v", u, reply.Message, strings.Join(reply.Details, ", "))
	}
	return &reply, nil
}

func Faucet(typ int) (*HelloBlockFaucetData, error) {
	u := fmt.Sprintf("%s?type=%d", helloBlockURL("v1/faucet"), typ)
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("error calling %v: %v", u, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var reply HelloBlockFaucetReply
	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("could not process JSON %v: %v", string(body), err)
	}
	if reply.Status != success {
		return nil, fmt.Errorf("error from %v: %v: %v", u, reply.Message, strings.Join(reply.Details, ", "))
	}
	return &reply.Data, nil
}

func Propagate(rawTxHex string) (*HelloBlockPropagateData, error) {
	u := helloBlockURL("v1/transactions")
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
		return nil, fmt.Errorf("could not process JSON %v: %v", string(body), err)
	}
	if reply.Status != success {
		return nil, fmt.Errorf("error from helloblock.io: %v: %v", reply.Message, strings.Join(reply.Details, ", "))
	}
	return &reply.Data, nil
}

func helloBlockURL(endpoint string) string {
	return fmt.Sprintf("https://%v.helloblock.io/%v", network.String(), endpoint)
}
