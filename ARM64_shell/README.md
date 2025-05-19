# AArch64 Shell

A minimal shell implementation written entirely in AArch64 assembly.
Demonstrates direct system call usage for basic shell functionality.

## Features

- Written in pure AArch64 assembly
- Uses Linux system calls only (no libc)
- Displays a prompt and reads user input
- Parses a command and one optional argument
- Forks a new process to execute commands
- Waits for child process to finish before continuing

## Assemble and Link

```bash
as shell.s -o shell.o
ld shell.o -o shell

## Example usage

Shell Initiated
$: /bin/ls
<output of ls>
$: /bin/echo hello
hello
$: exit
