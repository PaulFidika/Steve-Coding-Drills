#!/usr/bin/env -S uv run --script

# This shows how to run Claude Code programatically! This agent is programmable!

import subprocess

prompt = """

GIT checkout a NEW brnach.

Create ./cc_todo/todo.ts: a zero dependency library CLI todo app with basic CRUD.

THEN GIT stage, commit, and SWITCH back to main.

"""

command = ["claude", "-p", prompt, "--allowedTools", "Edit", "Bash", "Write"]

process = subprocess.run(command, check=True)

print("Claude process output: ", process.stdout)
