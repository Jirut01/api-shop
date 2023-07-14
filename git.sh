#!/bin/bash
if [ $# -eq 0 ]; then
        echo "No msg to commit"
        echo "Usage: " $0 " {msg commit}"
        echo "example: " $0 " update code"
        exit 1
else
    git add .
    git commit -m "$1"
    git push
fi
