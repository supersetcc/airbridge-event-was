#!/bin/bash

echo "[CodeDeploy] Stop supervisor processes..."
sudo service supervisor stop
echo "[CodeDeploy] OK"

echo "[CodeDeploy] Delete old source files..."
sudo rm -rf /var/www/airbridge-go-bypass-was/*
sudo rm -rf /var/www/airbridge-go-bypass-was/.*
echo "[CodeDeploy] OK"
