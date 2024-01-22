#!/bin/sh
set -ex
ipfs init
cp /mount/swarm.key ~/.ipfs/
ipfs bootstrap rm --all
ipfs bootstrap add /ip4/127.0.0.1/tcp/4001/p2p/12D3KooWQ3QnJLexZyEiTMdj7DjrXiT2HANLKLT2d4rB6LpwtuhP
mkdir /home/ipfsData 