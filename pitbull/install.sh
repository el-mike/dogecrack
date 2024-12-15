#!/bin/bash

# This script installs all necessary dependencies and does some basic setup.

apt -y install locales
apt -y install wget
apt -y install python3 python3-pip nano mc git python3-bsddb3
apt -y install libssl-dev build-essential automake pkg-config libtool libffi-dev libgmp-dev libyaml-cpp-dev libsecp256k1-dev
apt -y install tmux
sudo apt install python3-pyopencl
git clone https://github.com/3rdIteration/btcrecover.git
python3 -m pip install siphash24
pip3 install gdown
# Version needs to be specified at least for nvidia/opencl:runtime-ubuntu18.04 base image.
pip3 install pyopencl==2019.1.1
pip3 install -r ./btcrecover/requirements.txt
update-locale LANG=C.UTF-8
echo "set -g terminal-overrides \"xterm*:kLFT5=\eOD:kRIT5=\eOC:kUP5=\eOA:kDN5=\eOB:smkx@:rmkx@\"" > ~/.tmux.conf
