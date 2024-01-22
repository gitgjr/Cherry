#!/bin/sh
set -ex
apt install sudo
sudo apt upgrade
sudo apt -y install ffmpeg
cd /home
git clone https://github.com/gitgjr/Cherry.git
mkdir Cherry/data
mkdir Cherry/data/resource
mkdir Cherry/data/serverWork
mkdir Cherry/data/p2pWork

cp /resource/* /home/Cherry/data/resource

echo "Cherry Initialization Complete"