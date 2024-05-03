#!/bin/bash
until ./backend start; do
    echo "Program crashed with exit code $?.  Respawning.." >&2
    sleep 1
done