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
go mod init main
go mod tidy
go build -o gonotes_driver
cd ..
go mod init main
go mod tidy
go build -o gonotes
./gonotes