#!/bin/bash

# Attaches to tmux session running the Pitbull process.

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

source ./variables.sh

tmux a -t "$tmuxSessionName"
