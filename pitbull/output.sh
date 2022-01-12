#!/bin/bash

# Returns contents of progress_view.txt file.

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname


cat ./progress_view.txt
