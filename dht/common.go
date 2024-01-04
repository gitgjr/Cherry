package dht

import (
	"fmt"
	"main/zlog"
	"math/rand"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

func makePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		zlog.Error("Err on creating host: %v", zap.Error(err))
	}
	zlog.Debug("Parsed: ipfsAddr =", zap.Any("ipfsAddr", ipfsAddr))

	peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		zlog.Fatal("Err on creating peerIDStr: %v", zap.Error(err))
	}
	zlog.Debug("Parsed: PeerIDStr = ", zap.Any("peerIDStr", peerIDStr))

	peerID, err := peer.IDB58Decode(peerIDStr)
	if err != nil {
		zlog.Fatal("Err on decoding %s: %v", zap.Any("peerIDStr", peerIDStr), zap.Error(err))
	}
	zlog.Debug("Created peerID =", peerID)

	targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
	zlog.Debug("Created targetPeerAddr = ", zap.Any("targetPeerAddr", targetPeerAddr))

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	zlog.Debug("Decapsuated = ", zap.Any("targetAddr", targetAddr))

	return peerID, targetAddr
}

func generatePrivateKey(seed int64) crypto.PrivKey {
	randBytes := rand.New(rand.NewSource(seed))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)

	if err != nil {
		zlog.Error("DHT Could not generate Private Key", zap.Error(err))
	}

	return prvKey
}

type nullValidator struct{}

// Validate always returns success
func (nv nullValidator) Validate(key string, value []byte) error {
	zlog.Debug("DHT NullValidator Validate", zap.Any(key, value))
	return nil
}

// Select always selects the first record
func (nv nullValidator) Select(key string, values [][]byte) (int, error) {
	strs := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		strs[i] = string(values[i])
	}
	zlog.Debug("DHT NullValidator Select:", zap.Any(key, strs))

	return 0, nil
}
