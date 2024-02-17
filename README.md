<!-- SPDX-License-Identifier: CC-BY-4.0 -->

# `ghasum`

Checksums for GitHub Actions.

Compute and verify checksums for all GitHub Actions in a project to guarantee
that the Actions you choose to include haven't changed since. `ghasum` gives
better integrity guarantees than pinning Actions by commit hash and is also more
user friendly as well.

## Usage

To start using `ghasum` navigate to a project that use GitHub Actions and run:

```shell
ghasum init
```

Commit the `gha.sum` file that is created so that the checksums can be verified
in the future. To verify run:

```shell
ghasum verify
```

For further help with using `ghasum` run:

```shell
ghasum help
```

## Recommendations

When using ghasum it is recommend to pin all Actions to version tags. If Actions
are benign, these won't change over time. Major version tags or branch refs are
expected to change over time as changes are made to the Action, which results in
failing verification by ghasum. Commit SHAs do not have to be used because the
benefits they provide are covered by ghasum.

If an Action misbehaves - moving version refs after publishing - it is
recommended to use commit SHAs instead to avoid failing verification by ghasum.

```yaml
# Recommended: exact version tags
- uses: actions/checkout@v4.1.1

# Possible alternative: commit SHAs
- uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

# Discouraged: major version refs
- uses: actions/checkout@v4

# Discouraged: branches
- uses: actions/checkout@main
```

## Limitations

- Requires manual intervention when an Action is updated.
- The hashing algorithm used for checksums is not configurable.
- Checksums do not provide protections against [unpinnable actions].

[unpinnable actions]: https://www.paloaltonetworks.com/blog/prisma-cloud/unpinnable-actions-github-security/

## Background

The dependency ecosystem for GitHub Actions is fully reliant on git. The version
of an Action to use is specified using a git ref (branch or tag) or commit SHA.
Git refs provide no integrity guarantees. And while commit SHAs do provide some
integrity guarantees, since they're based on the older SHA1 hash the guarantees
are not optimal.

Besides being older and having better, modern algorithms available, SHA1 is
vulnerable to the [SHAttered] attack. This means it is possible for a motivated
and well-funded adversary to mount an attack on the Github Actions ecosystem.
GitHub does have [protections in place] to detect such an attack, but this is
specific to the [SHAttered] attack and, like hashing algorithms, probabilistic.

This project is a response to that theoretical attack - providing a way to get,
record, and validate checksums for GitHub Actions dependencies using a more
secure hashing algorithm. As an added benefit, it can also be used as an
alternative to in-workflow commit SHA.

[protections in place]: https://github.blog/2017-03-20-sha-1-collision-detection-on-github-com/
[shattered]: https://shattered.io/

## License

This software is available under the Apache License 2.0 license, see [LICENSE]
for the full license text. The contents of documentation are licensed under the
[CC BY 4.0] license.

[cc by 4.0]: https://creativecommons.org/licenses/by/4.0/
[LICENSE]: ./LICENSE
