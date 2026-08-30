#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <row> <column> <sprite_file>"
    exit 1
fi
tput sc
tput cup $1 $2
go run . -f $3 -sprite
tput rc

