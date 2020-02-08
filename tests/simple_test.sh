#!/bin/bash
LOG=./scheduled.log

# Cleanup
rm -f $LOG

# Start the daemon
./scheduled > $LOG 2>&1 &
sleep 2

# create a hostname-username.key file
./register
sleep 2

# Run a command
./schedule ls -l

sleep 30

# Stopping it 
pkill scheduled
