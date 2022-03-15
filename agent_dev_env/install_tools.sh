#!/bin/bash

echo "Installing Go...."
GO_VERSION=$(curl https://go.dev/VERSION?m=text)
curl -OL https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf ${GO_VERSION}.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" > ~/.profile
source ~/.profile

echo "Installing pip...."
curl -sSL https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python3 get-pip.py
echo "export PATH=$PATH:$(pwd)/.local/bin" >> ~/.profile
source ~/.profile

pip install invoke

echo "Installing cmake...."
sudo apt update
sudo apt install cmake -y
