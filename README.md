# `ghasum`

Checksums for GitHub Actions.

## Motivation

The dependency ecosystem for GitHub Actions is fully reliant on git. Dependencies can be specified
by a git ref (branch or tag) or commit SHA. Git refs provide no integrity guarantees. Commit SHAs do
provide integrity guarantees, but since they're based on SHA1 these guarantees are weaker then
necessary (see for example the [SHAttered] attack).

Hence, this project aims to provide a way to get, record, and validate checksums for GitHub Actions
dependencies using a more modern cryptographic hashing algorithm.

## TODO

An incomplete list of things that should be done:

- [ ] Better user experience
- [ ] Improve source code and use automated testing
- [ ] At least a basic continuous integration setup
- [ ] Record some sort of version in the checksum file in order to aid in changing the format or
      changing the hash algorithm in the future
- [ ] Convenient support to run as a GitHub Action to validate dependency integrity before the (rest
      of the) job runs

[shattered]: https://shattered.io/
