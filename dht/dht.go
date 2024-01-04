package dht

import (
	"context"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// var log = logging.Logger("main")

func addrForPort(p string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p))
}

func generateHost(ctx context.Context, port int64) (host.Host, *dht.IpfsDHT) {
	prvKey := generatePrivateKey(port)

	hostAddr, err := addrForPort(fmt.Sprintf("%d", port))
	if err != nil {
		fmt.Println(err)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(prvKey),
	}

	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		fmt.Println(err)
	}

	kadDHT, err := dht.New(ctx, host, dht.Validator(nullValidator{}))
	if err != nil {
		fmt.Println(err)
	}

	hostID := host.ID()
	fmt.Printf("Host MultiAddress: %s/ipfs/%s (%s)", host.Addrs()[0].String(), hostID.Pretty(), hostID.String())

	return host, kadDHT
}

func addPeers(ctx context.Context, h host.Host, kad *dht.IpfsDHT, peersArg string) {
	if len(peersArg) == 0 {
		return
	}

	peerStrs := strings.Split(peersArg, ",")
	for i := 0; i < len(peerStrs); i++ {
		peerID, peerAddr := makePeer(peerStrs[i])

		h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
		kad.Update(ctx, peerID)
	}
}
