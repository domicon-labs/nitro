package das

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"net/http"
	"testing"
	"time"
)

/*
*
curl -H "content-type: application/json" -X POST --data '{"id":0,"jsonrpc":"2.0","method":"eth_getFileDataByCommitment","params":["a595c834b09c0d6cbc37b6e6b9b7ac606ae32063630267a9151f1d695403020796e344441c051ddb6903fa4d713c283f"]}' http://13.212.115.195:8545

{"jsonrpc":"2.0","id":0,"result":{"sender":"0xf2fa2c2e6b3399237e0d5c413d21e37cf4db23b0","submitter":"0x1845b7295ae3ee0fc4b5fe60c05ea81637603764","length":"0x8","index":"0x1","commitment":"0xa595c834b09c0d6cbc37b6e6b9b7ac606ae32063630267a9151f1d695403020796e344441c051ddb6903fa4d713c283f","data":"0x0000000066274bf9","sign":"0xc96ae9f5069488d4db65449a5b28f098228abcfc17c2b000ec661b4bc889ef6a5d0d1600b1dd4ccc0ffdb0aaba9584745db9709026aad3a53675c1a03a17c3f01c","txhash":"0x36bd626b1ba3666ab19906fe30707ed6b5d91785099ce4c361b7892264fb427a"}}
*
*/
func TestDomiconGetFileDataByCommitment1(t *testing.T) {
	enableLogging()
	id := 0
	method := "eth_getFileDataByCommitment"
	commitment := "a595c834b09c0d6cbc37b6e6b9b7ac606ae32063630267a9151f1d695403020796e344441c051ddb6903fa4d713c283f"
	params := []string{commitment}
	peer := "http://13.212.115.195:8545"
	// 准备 JSON 数据
	jsonStr := []byte(fmt.Sprintf(`{"id":%d,"jsonrpc":"2.0","method":"%s","params":%q}`, id, method, params))

	// 发送 POST 请求
	resp, err := http.Post(peer, "application/json", bytes.NewBuffer(jsonStr))
	Require(t, err)

	// 读取响应内容
	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(resp.Body)
	Require(t, err)

	t.Log(respBody.String())
}

func TestDomiconGetFileDataByCommitment2(t *testing.T) {
	// 创建一个 RPC 客户端
	peer := "http://13.212.115.195:8545"
	method := "eth_getFileDataByCommitment"
	commitment := "a595c834b09c0d6cbc37b6e6b9b7ac606ae32063630267a9151f1d695403020796e344441c051ddb6903fa4d713c283f"

	client, err := rpc.DialHTTP(peer)
	Require(t, err)

	// 构建 JSON-RPC 请求体
	var result map[string]string
	err = client.CallContext(context.Background(), &result, method, commitment)
	Require(t, err)
	t.Log(result)
	t.Log(result["data"])
}

func TestDomiconStorageServiceGet(t *testing.T) {
	enableLogging()
	ctx := context.Background()
	commitment := "a595c834b09c0d6cbc37b6e6b9b7ac606ae32063630267a9151f1d695403020796e344441c051ddb6903fa4d713c283f"
	peer := "http://13.212.115.195:8545"

	svc, err := NewDomiconStorageService(ctx,
		DomiconStorageServiceConfig{
			Enable:      true,
			ReadTimeout: time.Minute,
			Peer:        peer,
		})
	Require(t, err)
	returnedData, err := svc.GetByCommitment(ctx, commitment)
	Require(t, err)
	t.Log(string(returnedData))
}
func TestConvert(t *testing.T) {
	data1 := "0000000066274bf9"
	t.Log(common.FromHex(data1))
	t.Log(common.Bytes2Hex(common.FromHex(data1)))

	data2 := "0x0000000066274bf9"
	t.Log(common.FromHex(data2))
	t.Log(common.Bytes2Hex(common.FromHex(data2)))
}
