#!/bin/bash

sudo docker stop ctr_neurone-ss
sudo docker rm ctr_neurone-ss
sudo docker rmi img_neurone-ss


sudo docker build --rm -t img_neurone-ss:latest .

sudo docker run --network="host" --name ctr_neurone-ss  -d img_neurone-ss
