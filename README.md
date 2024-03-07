# Distributed meeting recording storage

## How to run

Run `docker build imageName:imageVersion .` in ./docker/kubo0.25 to build image

Then run `docker run -t -d --name kubo-bootstrap  -v .\docker\kubo0.25\export:/export -v .\docker\kubo0.25\datas:/data/ipfs -p 1115:4001 -p 1115:4001/udp -p 1116:8080 -p 1117:5001 -p 1118:8081 golang:latest` to create a bootstrap node and run **bootstrapNodeinit.sh** to init the bootstrap node. After that copy Address of the bootstrap and add it to bootstrap list in other nodes.

Run `docker run -t -d -m 6144m --cpus=3 --name kubo-node3  -v .\docker\kubo0.25\mount:/mount -v .\data\resource:/resource -p 4001 -p 4001/udp -p 127.0.0.1:8082:8080 -p 127.0.0.1:5003:5001  gengjieran/go-ipfs:v2.1` to create a normal node. Run **nodeInit.sh** to init the kubo node and run **projectInit.sh** to init the project. 

Modify and run the main.go in cmd/Evaluate-p2p and cmd/Evaluate-server to do more test. Input videos need to have same resolution and frame rate.

You can see more commend in ./docker/dockerCommend.txt


