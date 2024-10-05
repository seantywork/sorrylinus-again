#!/bin/bash

# nginx

sudo apt update
sudo apt install nginx


# certbot

sudo snap install --classic certbot 

sudo ln -s /snap/bin/certbot /usr/bin/certbot 

# firewall


# sudo ufw allow 22

# sudo ufw allow 80

# sudo ufw allow 443

# sudo ufw allow 3000

# sudo ufw allow 3001
