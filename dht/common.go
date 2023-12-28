package dht

import (
	"fmt"
	"main/zlog"
	"math/rand"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

func makePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		zlog.Error("Err on creating host: %v", err)
	}
	zlog.Debug("Parsed: ipfsAddr = %s", ipfsAddr)

	peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalf("Err on creating peerIDStr: %v", err)
	}
	log.Debugf("Parsed: PeerIDStr = %s", peerIDStr)

	peerID, err := peer.IDB58Decode(peerIDStr)
	if err != nil {
		log.Fatalf("Err on decoding %s: %v", peerIDStr, err)
	}
	log.Debugf("Created peerID = %s", peerID)

	targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
	log.Debugf("Created targetPeerAddr = %v", targetPeerAddr)

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	log.Debugf("Decapsuated = %v", targetAddr)

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
