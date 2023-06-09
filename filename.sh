#!/bin/bash

CPU_NAME=$(lscpu | grep "Model name:" | sed -e 's/Model name://g'  | tr -d " :()@#+={};,.?[]")
NOW=$(date +"%Y-%m-%d_%H-%M-%S")
echo ${CPU_NAME}_${NOW}.out
