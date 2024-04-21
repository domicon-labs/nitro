package das

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/nitro/arbstate"
	"github.com/offchainlabs/nitro/das/dastree"
)

// DomiconStorageService implements DataAvailabilityReader
type DomiconStorageService struct {
	url string
}

func NewDomiconStorageService(protocol string, host string, port int) *DomiconStorageService {
	return &DomiconStorageService{
		url: fmt.Sprintf("%s://%s:%d", protocol, host, port),
	}
}

func NewDomiconStorageServiceFromURL(url string) (*DomiconStorageService, error) {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return nil, fmt.Errorf("protocol prefix 'http://' or 'https://' must be specified for DomiconStorageService; got '%s'", url)

	}
	return &DomiconStorageService{
		url: url,
	}, nil
}

func (c *DomiconStorageService) GetByHash(ctx context.Context, hash common.Hash) ([]byte, error) {
	res, err := http.Get(c.url + getByHashRequestPath + EncodeStorageServiceKey(hash))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response RestfulDasServerResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(response.Data)))
	decodedBytes, err := io.ReadAll(decoder)
	if err != nil {
		return nil, err
	}
	if !dastree.ValidHash(hash, decodedBytes) {
		return nil, arbstate.ErrHashMismatch
	}

	return decodedBytes, nil
}

func (c *DomiconStorageService) HealthCheck(ctx context.Context) error {
	res, err := http.Get(c.url + healthRequestPath)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
	return nil
}

func (c *DomiconStorageService) ExpirationPolicy(ctx context.Context) (arbstate.ExpirationPolicy, error) {
	res, err := http.Get(c.url + expirationPolicyRequestPath)
	if err != nil {
		return -1, err
	}
	if res.StatusCode != http.StatusOK {
		return -1, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return -1, fmt.Errorf("HTTP error with status %d returned by server: %s", res.StatusCode, http.StatusText(res.StatusCode))
	}

	var response RestfulDasServerResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	return arbstate.StringToExpirationPolicy(response.ExpirationPolicy)
}
