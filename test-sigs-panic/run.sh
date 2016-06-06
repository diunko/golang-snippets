#!/bin/bash


trap "echo INT" SIGINT
trap "echo TERM" SIGTERM

while true; do
  echo "child's here!"
  sleep 1
done
