package dht

import (
	"context"
	"log"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func NewKDHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (dis, error) {
	var options []dht.Option

	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kademliaDHT, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return NewRoutingDiscovery(kademliaDHT), nil
}
