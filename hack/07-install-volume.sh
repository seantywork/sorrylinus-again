#!/bin/bash


set -euxo pipefail

#######

# nfs

local_ip=$1

apt-get update

apt-get install -y nfs-kernel-server

apt-get install -y nfs-common

mkdir -p /data

chown -R root:root /data

chmod -R 777 /data



IP=$local_ip

cat << EOF > /etc/exports

/data $IP/24(rw,sync,no_subtree_check,no_root_squash) 10.10.0.0/16(rw,sync,no_subtree_check,no_root_squash)

EOF

exportfs -a

systemctl restart nfs-kernel-server


# nfs provisioner

set -euxo pipefail

#######

server_ip=$1


sudo helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/


sudo helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner --set nfs.server=$server_ip --set nfs.path=/data


sudo kubectl apply -f ../k8s/default-storage-class.yaml


######


