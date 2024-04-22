package das

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/offchainlabs/nitro/arbstate"
	"github.com/offchainlabs/nitro/util/pretty"
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
	log.Trace("das.DomiconStorageService.GetByHash", "hash", pretty.PrettyHash(hash))

	id := 0
	method := "eth_getFileDataByCommitment"
	params := []string{hash.String()}

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

//	func NewDomiconStorageServiceFromURL(url string) (*DomiconStorageService, error) {
//		if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
//			return nil, fmt.Errorf("protocol prefix 'http://' or 'https://' must be specified for DomiconStorageService; got '%s'", url)
//
//		}
//		return &DomiconStorageService{
//			url: url,
//		}, nil
//	}
//
//	func (c *DomiconStorageService) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
//		res, err := http.Get(c.url + getByHashRequestPath + EncodeStorageServiceKey(hash))
//		if err != nil {
//			return nil, err
//		}
//		if res.StatusCode != http.StatusOK {
//			return nil, fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
//		}
//
//		body, err := io.ReadAll(res.Body)
//		if err != nil {
//			return nil, err
//		}
//
//		var response RestfulDasServerResponse
//		err = json.Unmarshal(body, &response)
//		if err != nil {
//			return nil, err
//		}
//
//		decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(response.Data)))
//		decodedBytes, err := io.ReadAll(decoder)
//		if err != nil {
//			return nil, err
//		}
//		if !dastree.ValidHash(hash, decodedBytes) {
//			return nil, arbstate.ErrHashMismatch
//		}
//
//		return decodedBytes, nil
//	}
func (c *DomiconStorageService) HealthCheck(ctx context.Context) error {
	//res, err := http.Get(c.url + healthRequestPath)
	//if err != nil {
	//	return err
	//}
	//if res.StatusCode != http.StatusOK {
	//	return fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
	//}
	return nil
}

//
//func (c *DomiconStorageService) ExpirationPolicy(ctx context.Context) (arbstate.ExpirationPolicy, error) {
//	res, err := http.Get(c.url + expirationPolicyRequestPath)
//	if err != nil {
//		return -1, err
//	}
//	if res.StatusCode != http.StatusOK {
//		return -1, err
//	}
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		return -1, fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
//	}
//
//	var response RestfulDasServerResponse
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		return -1, err
//	}
//
//	return arbstate.StringToExpirationPolicy(response.ExpirationPolicy)
//}
