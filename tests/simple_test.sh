#!/bin/bash
LOG=./scheduled.log

# Cleanup
rm -f $LOG
#rm -rf ./data

# Start the daemon
./scheduled > $LOG 2>&1 &
sleep 2

# create a hostname-username.key file
./register
sleep 2

# Run a command TODO: currently needs full path :-(
./schedule /bin/ps -elf

sleep 5 

# Stopping it 
pkill scheduled
