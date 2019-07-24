package bitshares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/blocktree/bitshares-adapter/types"
	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

// WalletClient is a Bitshares RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type WalletClient struct {
	WalletAPI, ServerAPI string
	Debug                bool
	client               *req.Req
}

// NewWalletClient init a rpc client
func NewWalletClient(serverAPI, walletAPI string, debug bool) *WalletClient {

	walletAPI = strings.TrimSuffix(walletAPI, "/")
	serverAPI = strings.TrimSuffix(serverAPI, "/")
	c := WalletClient{
		WalletAPI: walletAPI,
		ServerAPI: serverAPI,
		Debug:     debug,
	}

	api := req.New()
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *WalletClient) call(method string, request interface{}, queryWalletAPI bool) (*gjson.Result, error) {

	var (
		body = make(map[string]interface{}, 0)
	)

	if c.client == nil {
		return nil, fmt.Errorf("API url is not setup. ")
	}

	authHeader := req.Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	//json-rpc
	body["jsonrpc"] = "2.0"
	body["id"] = 1
	body["method"] = method
	body["params"] = request

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	host := c.ServerAPI
	if queryWalletAPI {
		host = c.WalletAPI
	}

	r, err := c.client.Post(host, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = c.isError(r)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

// isError 是否报错
func (c *WalletClient) isError(r *req.Resp) error {

	if r.Response().StatusCode != http.StatusOK {
		message := r.Response().Status
		status := r.Response().StatusCode
		return fmt.Errorf("[%d]%s", status, message)
	}

	result := gjson.ParseBytes(r.Bytes())

	if result.Get("error").IsObject() {

		return fmt.Errorf("[%d]%s",
			result.Get("error.code").Int(),
			result.Get("error.message").String())

	}

	return nil

}

// GetObjects return a block by the given block number
func (c *WalletClient) GetObjects(assets ...types.ObjectID) (*gjson.Result, error) {
	resp, err := c.call("get_objects", []interface{}{objectsToParams(assets)}, false)
	return resp, err
}

func objectsToParams(objs []types.ObjectID) []string {
	objsStr := make([]string, len(objs))
	for i, asset := range objs {
		objsStr[i] = asset.String()
	}
	return objsStr
}

// GetBlockchainInfo returns current blockchain data
func (c *WalletClient) GetBlockchainInfo() (*BlockchainInfo, error) {
	r, err := c.call("get_dynamic_global_properties", []interface{}{}, false)
	if err != nil {
		return nil, err
	}
	info := NewBlockchainInfo(r)
	return info, nil
}

// GetBlockByHeight returns a certain block
func (c *WalletClient) GetBlockByHeight(height uint32) (*Block, error) {
	r, err := c.call("get_block_header", []interface{}{height + 1}, false)
	if err != nil {
		return nil, err
	}
	header := NewBlockHeader(r)

	r, err = c.call("get_block", []interface{}{height}, false)
	if err != nil {
		return nil, err
	}
	block := NewBlock(height, r)
	block.BlockID = header.Previous

	// block.CalculateID()
	// log.Std.Info("calculated block id:%s\n", block.BlockID)

	return block, nil
}

// GetTransaction returns the TX
func (c *WalletClient) GetTransaction(height uint32, trxInBlock int) (*types.Transaction, error) {
	r, err := c.call("get_transaction", []interface{}{height, trxInBlock}, false)
	if err != nil {
		return nil, err
	}
	if r.Raw == "null" {
		return nil, fmt.Errorf("cannot find this transaction: %v, %v", height, trxInBlock)
	}
	block, err := c.GetBlockByHeight(height)
	if err != nil {
		return nil, err
	}
	if len(block.TransactionIDs) <= trxInBlock {
		return nil, fmt.Errorf("cannot find this transaction on the block: %v, %v", height, trxInBlock)
	}
	return NewTransaction(r, block.TransactionIDs[trxInBlock])
}

// GetAssetsBalance Returns information about the given account.
func (c *WalletClient) GetAssetsBalance(account types.ObjectID, asset types.ObjectID) (*Balance, error) {
	r, err := c.call("get_account_balances", []interface{}{account.String(), []interface{}{asset.String()}}, false)
	if err != nil {
		return nil, err
	}
	return NewBalance(r), nil
}

// GetAssetsBalance Returns information about the given account.
func (c *WalletClient) GetAccountID(name string) (*types.ObjectID, error) {
	r, err := c.call("lookup_accounts", []interface{}{name, 1}, false)
	if err != nil {
		return nil, err
	}
	arr := r.Array()
	if len(arr) > 0 {
		if arr[0].Array()[0].String() == name {
			id := arr[0].Array()[1].String()
			objectID := types.MustParseObjectID(id)
			return &objectID, nil
		}
	}
	return nil, fmt.Errorf("[%s] have not registered", name)
}

// GetAssetsBalance Returns information about the given account.
func (c *WalletClient) GetAccounts(names_or_ids []string) ([]*types.Account, error) {
	var resp []*types.Account
	r, err := c.call("get_accounts", []interface{}{names_or_ids}, false)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(r.Raw), &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *WalletClient) GetRequiredFee(ops []types.Operation, assetID string) ([]types.AssetAmount, error) {
	var resp []types.AssetAmount

	opsJSON := []interface{}{}
	for _, o := range ops {
		_, err := json.Marshal(o)
		if err != nil {
			return []types.AssetAmount{}, err
		}

		opArr := []interface{}{o.Type(), o}

		opsJSON = append(opsJSON, opArr)
	}
	r, err := c.call("get_required_fees", []interface{}{opsJSON, assetID}, false)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(r.Raw), resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// BroadcastTransaction broadcast a transaction
func (c *WalletClient) BroadcastTransaction(tx *types.Transaction) (*BroadcastResponse, error) {
	resp := BroadcastResponse{}

	r, err := c.call("broadcast_transaction", []interface{}{tx}, false)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(r.Raw), resp); err != nil {
		return nil, err
	}
	return &resp, err
}

// GetTransactionID return the TX ID
func (c *WalletClient) GetTransactionID(tx *types.Transaction) (string, error) {
	r, err := c.call("get_transaction_id", []interface{}{tx}, true)
	if err != nil {
		return "", err
	}
	return r.String(), err
}
