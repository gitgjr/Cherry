#!/bin/sh
set -ex
wget https://dist.ipfs.tech/kubo/v0.25.0/kubo_v0.25.0_linux-amd64.tar.gz
tar xvfz kubo_v0.25.0_linux-amd64.tar.gz
./kubo/install.sh
rm kubo_v0.25.0_linux-amd64.tar.gz
rm -R kubo

ipfs init
echo -e "/key/swarm/psk/1.0.0/\n/base16/\n`tr -dc 'a-f0-9' < /dev/urandom | head -c64`" > ~/.ipfs/swarm.key
ipfs bootstrap rm --all
