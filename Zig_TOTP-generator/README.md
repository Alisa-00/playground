# TOTP generator

A simple command-line TOTP (Time-based One-Time Password) generator written in [Zig](https://ziglang.org/).

This tool lets you store shared secrets and generate 6-digit time-based codes compatible with 2 factor authentication services like Google Authenticator. It is based on [RFC 6238](https://datatracker.ietf.org/doc/html/rfc6238).

## Features

- Add and store TOTP secrets for multiple accounts
- Generate current TOTP codes
- List stored accounts
- Stores secrets locally in a simple JSON file
- Base32 decoding of secrets (compatible with most authenticator apps)

## Usage

### Add a new account
```bash
totp add github JBSWY3DPEHPK3PXP
Added github to vault successfully!
```

### Get current TOTP code
```sh
totp get github
github: 123456
```

### List stored accounts
```sh
totp list
github: 123456
account: 654321
...
```

## Notes
* Secrets are stored in a local file (e.g., ~/.vaultfile).
* Secrets must be Base32-encoded
* TOTP codes are generated using HMAC-SHA1
* This tool was created for learning purposes only. It does not encrypt stored secrets and so it is NOT SECURE to use.
