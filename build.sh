#!/bin/bash
DIR=$1
echo "gom build -o bin/cuttle cuttle/*"
gom build -o bin/cuttle cuttle/*
