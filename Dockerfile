FROM golang:latest

WORKDIR /

# Generic
ADD config/config.yaml config/config.yaml
ADD api api


