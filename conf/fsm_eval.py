#!/usr/bin/python3

import sqlite3
import sys
from dataclasses import dataclass


ARGV = sys.argv

@dataclass
class Args:
    entry: str
    type: str
    id: str
    leader_id: str
    role: str
    last_log_index: str
    current_term: str
    voted_for: str

args = Args(
    entry = ARGV[1],
    type = ARGV[2],
    id = ARGV[3],
    leader_id = ARGV[4],
    role = ARGV[5],
    last_log_index = ARGV[6],
    current_term = ARGV[7],
    voted_for = ARGV[8],
)

con = sqlite3.connect(f"conf/{args.id}.db")
cur = con.cursor()

try:
    if args.type == "COMMAND":
        cur.executescript(args.entry)
        con.commit()

    if args.type == "QUERY":
        for row in cur.execute(args.entry):
            print(row)
except Exception as e:
    sys.exit(str(e))
finally:
    con.close()
