package das

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/offchainlabs/nitro/arbstate"
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
	Enable:      true,
	ReadTimeout: time.Minute,
	Peer:        "http://13.212.115.195:8545",
}

func DomiconStorageServiceConfigAddOptions(prefix string, f *flag.FlagSet) {
	f.Bool(prefix+".enable", DefaultDomiconStorageServiceConfig.Enable, "enable storage/retrieval of sequencer batch data from domicon")
	f.Duration(prefix+".read-timeout", DefaultDomiconStorageServiceConfig.ReadTimeout, "timeout for domicon reads, since by default it will wait forever. Treat timeout as not found")
	f.String(prefix+".peer", DefaultDomiconStorageServiceConfig.Peer, "domicon peer to connect to, eg http://xx.xx.xx:8545")
}

type DomiconStorageService struct {
	config        DomiconStorageServiceConfig
	domiconClient *rpc.Client
}

func NewDomiconStorageService(ctx context.Context, config DomiconStorageServiceConfig) (*DomiconStorageService, error) {
	client, err := rpc.DialHTTP(config.Peer)
	if err != nil {
		log.Warn("failed to dial rpc endpoint: %v, err: %v", config.Peer, err)
		return nil, err
	}

	return &DomiconStorageService{
		config:        config,
		domiconClient: client,
	}, nil
}

func (s *DomiconStorageService) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	log.Info("das.DomiconStorageService.GetByHash", "hash", pretty.PrettyHash(hash))

	return nil, nil
}

func (s *DomiconStorageService) GetByCommitment(ctx context.Context, commitment string) ([]byte, error) {
	log.Info("das.DomiconStorageService.GetByCommitment", "commitment", commitment)

	method := "eth_getFileDataByCommitment"

	var result map[string]string
	err := s.domiconClient.CallContext(ctx, &result, method, commitment)
	if err != nil {
		log.Warn("Error sending request:", err)
		return nil, err
	}

	data := common.FromHex(result["data"])

	return data, nil
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
