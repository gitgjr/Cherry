package ipfs

import (
	"main/zlog"

	"github.com/ipfs/kubo/client/rpc"
	"go.uber.org/zap"
)

// type IPFSClient int

func NewClient() {
	node, err := rpc.NewLocalApi()
	if err != nil {
		zlog.Fatal("can not connect to IPFS client", zap.Error(err))
	}
	return node
}
