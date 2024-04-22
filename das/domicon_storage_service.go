package das

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/offchainlabs/nitro/arbstate"
	"net/http"
	"time"

	flag "github.com/spf13/pflag"
)

type DomiconStorageServiceConfig struct {
	Enable      bool          `koanf:"enable"`
	ReadTimeout time.Duration `koanf:"read-timeout"`
	Peer        string        `koanf:"peer"`
}

var DefaultDomiconStorageServiceConfig = DomiconStorageServiceConfig{
	Enable:      false,
	ReadTimeout: time.Minute,
	Peer:        "",
}

func DomiconStorageServiceConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.Bool(prefix+".enable", DefaultDomiconStorageServiceConfig.Enable, "enable storage/retrieval of sequencer batch data from domicon")
	f.Duration(prefix+".read-timeout", DefaultDomiconStorageServiceConfig.ReadTimeout, "timeout for domicon reads, since by default it will wait forever. Treat timeout as not found")
	f.String(prefix+".peer", DefaultDomiconStorageServiceConfig.Peer, "domicon peer to connect to, eg http://xx.xx.xx:8545")
}

type DomiconStorageService struct {
	config DomiconStorageServiceConfig
}

func NewDomiconStorageService(ctx context.Context, config DomiconStorageServiceConfig) (*DomiconStorageService, error) {
	return &DomiconStorageService{
		config: config,
	}, nil
}

func (s *DomiconStorageService) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	return nil, nil
}

func (s *DomiconStorageService) GetByCommitment(ctx context.Context, commitment []byte) ([]byte, error) {
	log.Trace("das.DomiconStorageService.GetByCommitment", "commitment", hex.EncodeToString(commitment))

	id := 0
	method := "eth_getFileDataByCommitment"
	params := []string{hex.EncodeToString(commitment)}

	// 准备 JSON 数据
	jsonStr := []byte(fmt.Sprintf(`{"id":%d,"jsonrpc":"2.0","method":"%s","params":%q}`, id, method, params))

	// 发送 POST 请求
	resp, err := http.Post(s.config.Peer, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Warn("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		log.Warn("Error reading response:", err)
		return nil, err
	}

	log.Trace("Response:", respBody.String())

	return respBody.Bytes(), nil
}

func (s *DomiconStorageService) Put(ctx context.Context, data []byte, timeout uint64) error {
	return nil
}

func (s *DomiconStorageService) ExpirationPolicy(ctx context.Context) (arbstate.ExpirationPolicy, error) {
	return arbstate.KeepForever, nil
}

func (s *DomiconStorageService) Sync(ctx context.Context) error {
	return nil
}

func (s *DomiconStorageService) Close(ctx context.Context) error {
	return nil
}

func (s *DomiconStorageService) String() string {
	return "DomiconStorageService"
}

func (s *DomiconStorageService) HealthCheck(ctx context.Context) error {
	return nil
}
