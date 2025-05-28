# Base64 Encoder/Decoder in Zig

A simple command-line tool for Base64 encoding and decoding, written in [Zig](https://ziglang.org/).
It takes input from the command line and encodes it into Base64 or decodes Base64 back into raw data.

## Features

- Encode strings into Base64
- Decode Base64 strings (WIP)
- Minimal, dependency-free implementation using Zig standard library

## Usage

```sh
zig run -- [encode|decode] [text]
```

or
```sh
zig build-exe main.zig
./main [encode|decode] [text]
```

## Examples

### Encode:

```sh
zig run -- encode "hello"
aGVsbG8=
```

### Decode (WIP):

```sh
zig run -- decode "aGVsbG8="
hello
```
