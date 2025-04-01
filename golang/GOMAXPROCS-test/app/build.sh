#!/bin/bash

VERSION=0.2

docker build -t mansoor1/golang-bench:$VERSION . 
docker push mansoor1/golang-bench:$VERSION