#!/bin/bash

# For Vast.ai images, we need to run this as a part of onstart.sh script, otherwise the PATH set up in
# Dockerfile may be overridden by some other launch operations.
echo "export PATH=/app:${PATH}" >> /etc/environment
