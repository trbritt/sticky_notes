#!/usr/bin/bash
# .PHONY: driver all main

# all: driver main 
# ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# driver:
# 	$(CD) /home/trbritt/Desktop$(ROOT_DIR)/driver 
# 	go build -o gonotes_driver
# 	$(CD) ..

# main:
# 	go build -o gonotes 

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