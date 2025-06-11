#!/usr/bin/env bash

sketch=$(ls | grep *.ino)

fqbn="arduino:avr:nano:cpu=atmega328old"
device="/dev/ttyUSB0"

arduino-cli compile -v --fqbn $fqbn "./$sketch"

arduino-cli upload -p $device -v --fqbn $fqbn "./$sketch"