#!/bin/bash

grep -q BCM /proc/cpuinfo || ( echo -e "\x1b[1;31mPlease run the installer on a Raspberry pi\x1b[1;0m" && exit 1 )

echo -e "\x1b[1;34minstallation running\x1b[1;0m"

sudo mkdir -p /usr/bin/smarthome-hw-ir || exit 1
sudo chown -R pi /usr/bin/smarthome-hw-ir  || exit 1
mv ./smarthome-hw-ir /usr/bin/smarthome-hw-ir/ || exit 1
sudo cp ./smarthome-hw-ir.service /lib/systemd/system/smarthome-hw-ir.service || exit 1

# Reload systemd
sudo systemctl daemon-reload || exit 1
sudo systemctl enable smarthome-hw-ir --now || exit 1

echo -e "\x1b[1;32minstallation completed\x1b[1;0m"
