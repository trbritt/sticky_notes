#!/usr/bin/bash

cd /home/trbritt/Desktop/projects/sticky_notes/driver
if [ ! -f go.mod ]; then
    go mod init main
    go mod tidy
fi  
if [ ! -f gonotes_driver ]; then
    go build -o gonotes_driver
fi 
cd ..
if [ ! -f go.mod ]; then
    go mod init main
    go mod tidy
fi
if [ ! -f gonotes ]; then
    go build -o gonotes
fi
./gonotes