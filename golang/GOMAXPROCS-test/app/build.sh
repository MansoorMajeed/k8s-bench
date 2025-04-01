#!/bin/bash

docker build -t golang-bench:$(git rev-parse --short HEAD) . 