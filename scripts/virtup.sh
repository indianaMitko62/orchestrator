#!/bin/bash

qemu-system-x86_64 -display none -enable-kvm -cpu host -smp 16 -boot menu=on -drive file=/var/lib/libvirt/images/centos7.0.qcow2,if=virtio,id=root -m 16G -net bridge,br=br-orch -name node1 &
qemu-system-x86_64 -display none -enable-kvm -cpu host -smp 16 -boot menu=on -drive file=/var/lib/libvirt/images/centos7.0-clone.qcow2,if=virtio,id=root -m 16G -net nic,model=virtio -net bridge,br=br-orch -name node2 &
qemu-system-x86_64 -display none -enable-kvm -cpu host -smp 16 -boot menu=on -drive file=/var/lib/libvirt/images/centos7.0-clone-clone.qcow2,if=virtio,id=root -m 16G -net nic,model=virtio -net bridge,br=br-orch -name node3 &

wait
wait
wait
