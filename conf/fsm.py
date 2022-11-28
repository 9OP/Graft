#!/usr/bin/python3

import sqlite3
import sys
import json
import os
import dataclasses
import enum

class EvalType(enum.IntEnum):
    INIT = 1
    QUERY = 2
    COMMAND = 3

@dataclasses.dataclass
class Input:
    id: str
    type: int
    data: str

def get_input() -> Input:
    payload = json.loads(sys.argv[1])
    return Input(
        id=payload["id"],
        type=payload["type"],
        data=payload["data"],
    )

def init(input: Input):
    os.remove(f"conf/{input.id}.db")

def command(input: Input):
    con = sqlite3.connect(f"conf/{input.id}.db")
    cur = con.cursor()
    cur.executescript(input.data)
    con.commit()
    con.close()

def query(input: Input):
    # Copy db in memory to prevent writes
    con = sqlite3.connect(":memory:")
    source = sqlite3.connect(f"conf/{input.id}.db")
    source.backup(con)
    con.row_factory = sqlite3.Row
    cur = con.cursor()
    res = []
    for row in cur.execute(input.data):
        res.append(dict(row))
    print(json.dumps(res, indent=2))

if __name__ == "__main__":
    input = get_input()
    if input.type == EvalType.INIT:
        init(input)
    elif input.type == EvalType.QUERY:
        query(input)
    elif input.type == EvalType.COMMAND:
        command(input)
