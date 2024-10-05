#!/bin/bash


sudo mkdir -p /data/sorrylinus-again


sudo mkdir -p /data/sorrylinus-again/admin

sudo echo '{}' | sudo tee /data/sorrylinus-again/admin/seantywork@gmail.com.json

sudo mkdir -p /data/sorrylinus-again/media

sudo mkdir -p /data/sorrylinus-again/media/article

sudo mkdir -p /data/sorrylinus-again/media/image

sudo mkdir -p /data/sorrylinus-again/media/video

sudo mkdir -p /data/sorrylinus-again/session


sudo mkdir -p /data/sorrylinus-again/user

sudo mkdir -p /data/sorrylinus-again/log

sudo echo '{}' | sudo tee /data/sorrylinus-again/log/stat.json



sudo chmod -R 777 /data/sorrylinus-again

sudo kubectl apply -n frank -f ../k8s/app-volume.yaml