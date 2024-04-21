package das

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/offchainlabs/nitro/util/pretty"
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

	return nil, nil

	//oracle := func(h common.Hash) ([]byte, error) {
	//	thisCid, err := hashToCid(h)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	ipfsPath := path.IpfsPath(thisCid)
	//	log.Trace("Retrieving IPFS path", "path", ipfsPath.String())
	//
	//	parentCtx := ctx
	//	if doPin {
	//		// If we want to pin this batch, then detach from the parent context so
	//		// we are not canceled before s.config.ReadTimeout.
	//		parentCtx = context.Background()
	//	}
	//
	//	timeoutCtx, cancel := context.WithTimeout(parentCtx, s.config.ReadTimeout)
	//	defer cancel()
	//	rdr, err := s.ipfsApi.Block().Get(timeoutCtx, ipfsPath)
	//	if err != nil {
	//		if timeoutCtx.Err() != nil {
	//			return nil, ErrNotFound
	//		}
	//		return nil, err
	//	}
	//
	//	data, err := io.ReadAll(rdr)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	if doPin {
	//		go func() {
	//			pinCtx, pinCancel := context.WithTimeout(context.Background(), s.config.ReadTimeout)
	//			defer pinCancel()
	//			err := s.ipfsApi.Pin().Add(pinCtx, ipfsPath)
	//			// Recursive pinning not needed, each dastree preimage fits in a single
	//			// IPFS block.
	//			if err != nil {
	//				// Pinning is best-effort.
	//				log.Warn("Failed to pin in IPFS", "hash", pretty.PrettyHash(hash), "path", ipfsPath.String())
	//			} else {
	//				log.Trace("Pin in IPFS successful", "hash", pretty.PrettyHash(hash), "path", ipfsPath.String())
	//			}
	//		}()
	//	}
	//
	//	return data, nil
	//}

	//return dastree.Content(hash, oracle)
}

//func NewDomiconStorageServiceFromURL(url string) (*DomiconStorageService, error) {
//	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
//		return nil, fmt.Errorf("protocol prefix 'http://' or 'https://' must be specified for DomiconStorageService; got '%s'", url)
//
//	}
//	return &DomiconStorageService{
//		url: url,
//	}, nil
//}
//
//func (c *DomiconStorageService) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
//	res, err := http.Get(c.url + getByHashRequestPath + EncodeStorageServiceKey(hash))
//	if err != nil {
//		return nil, err
//	}
//	if res.StatusCode != http.StatusOK {
//		return nil, fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
//	}
//
//	body, err := io.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var response RestfulDasServerResponse
//	err = json.Unmarshal(body, &response)
//	if err != nil {
//		return nil, err
//	}
//
//	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(response.Data)))
//	decodedBytes, err := io.ReadAll(decoder)
//	if err != nil {
//		return nil, err
//	}
//	if !dastree.ValidHash(hash, decodedBytes) {
//		return nil, arbstate.ErrHashMismatch
//	}
//
//	return decodedBytes, nil
//}
//
//func (c *DomiconStorageService) HealthCheck(ctx context.Context) error {
//	res, err := http.Get(c.url + healthRequestPath)
//	if err != nil {
//		return err
//	}
//	if res.StatusCode != http.StatusOK {
//		return fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
//	}
//	return nil
//}
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
