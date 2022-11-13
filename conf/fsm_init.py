#!/usr/bin/python3

import sys
import os

ARGV = sys.argv
id = ARGV[1]
path = f"conf/{id}.db"
print("remove", path)
os.remove(path)
