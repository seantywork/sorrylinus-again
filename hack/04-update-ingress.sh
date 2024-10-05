#!/bin/bash




/bin/cp -Rf ../nginx/fd.conf /etc/nginx/conf.d/fd.conf 

/bin/cp -Rf ../nginx/nginx.conf /etc/nginx/nginx.conf

systemctl restart nginx

certbot --nginx --rsa-key-size 4096 --no-redirect 


systemctl restart nginx
