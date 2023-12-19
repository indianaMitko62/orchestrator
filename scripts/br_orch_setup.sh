#!/bin/bash

sudo ip link add br-orch type bridge; sudo ip link set br-orch up

sudo ip tuntap add dev nic-orch mode tap

sudo ip link set nic-orch up promisc on

sudo ip link set nic-orch master br-orch

sudo ip addr add 192.168.42.1/24 broadcast 192.168.42.255 dev br-orch
