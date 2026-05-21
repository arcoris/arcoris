# Security: arcoris.dev/testutil

## Scope

This module contains test-only helpers.

## Security model

- no production runtime behavior;
- no authentication;
- no authorization;
- no cryptography;
- no sandboxing;
- no isolation.

## Misuse

- do not import from production code;
- do not treat panic assertions as recovery or safety mechanisms.

## Reporting

Security issues should be reported through the main ARCORIS security process.
