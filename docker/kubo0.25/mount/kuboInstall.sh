#!/bin/sh
set -ex
wget https://dist.ipfs.tech/kubo/v0.25.0/kubo_v0.25.0_linux-amd64.tar.gz
tar xvfz kubo_v0.25.0_linux-amd64.tar.gz
./kubo/install.sh
rm kubo_v0.25.0_linux-amd64.tar.gz
rm -R kubo

