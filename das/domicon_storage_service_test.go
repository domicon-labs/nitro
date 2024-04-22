package das

import (
	"context"
	"testing"
	"time"
)

func TestDomiconStorageServiceGet(t *testing.T) {
	enableLogging()
	ctx := context.Background()
	svc, err := NewDomiconStorageService(ctx,
		DomiconStorageServiceConfig{
			Enable:      true,
			ReadTimeout: time.Minute,
			Peer:        "http://54.177.136.244:8545",
		})
	defer svc.Close(ctx)
	Require(t, err)
	//returnedData, err := svc.GetByHash(ctx, common.BytesToHash(hash))
	//Require(t, err)
	//if !bytes.Equal(data, returnedData) {
	//	Fail(t, "Returned data didn't match!")
	//}
}
