#!/bin/bash

echo "Installing Go"
GO_VERSION=$(curl https://go.dev/VERSION?m=text)
curl -OL https://golang.org/dl/${GO_VERSION}.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf ${GO_VERSION}.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile

echo "Installing fzf...."
curl -OL https://github.com/junegunn/fzf/releases/download/0.29.0/fzf-0.29.0-linux_amd64.tar.gz
sudo tar -C /usr/bin -xvf fzf-0.29.0-linux_amd64.tar.gz

echo "Installing pip"
curl -sSL https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python3 get-pip.py
echo "export PATH=$PATH:$(pwd)/.local/bin" >> ~/.profile
source ~/.profile

pip install invoke

echo "Installing cmake"
pip install cmake --upgrade


echo "Installing .gitconfig"
curl -OL https://raw.githubusercontent.com/keisku/dotfiles/main/.gitconfig
