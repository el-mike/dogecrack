#!/bin/bash

# Stores common variables, that are reused between scripts.

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

viewFile='progress_view.txt'
defaultPasslistFile='passlist.txt'
errLogFile='err_log.txt'

pipe='btcrecover_out'

tmuxSessionName='pitbull'
