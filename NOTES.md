# SHAttered git commits

## Overview

1. The [SHAttered website](https://shattered.io/) states it's possible to make malicious git commits with the same SHA as an existing commit.
2. However, "the attack required 'the equivalent processing power of 6,500 years of single-CPU computations {...}.'" So it's not all that feasible in practice
3. From [SHA1 on Wikipedia](https://en.wikipedia.org/wiki/SHA-1), a more recent attack requires approximately the same amount of computation "{...} at the time of publication would cost US$45K per generated collision.

## Attack Description (in detail)

### Steps

1. Create valid commit and advertise it as a release (e.g. with a git tag or GitHub release).
2. Get people to use it.
3. Craft a malicious commit with the same hash.
4. People using that release now use malicious code which they can't tell from the integrity hash.

### Discussion

- I don't think that the older the release is (in terms of descendant commits) the harder this attack becomes. Namely, the child of a commit's only "chaining" to the parent is the parent's commit hash. So, since the commit hash didn't change there's no need to find a collision for more than one commit. The implication is that this attack can be carried out for any historical commit.
- Open questions:
  - Can existing clones be used to detect this attack? What is git's behavior locally when the hash of a historic commit is unchanged but the contents did change?
  - Can forks be used to detect this attack?

### Impact

- GitHub Actions
  - Dependencies in GitHub Actions may be specified/identified/versioned by the commit SHA of the git repository that hosts the dependency (referred to as an "Action"). This is the most secure approach, alternatively it can be specified/identified/versioned using git refs (such as branches or tags) but those are inherently mutable.
  - Hence, as a result of the attack GitHub Actions users that use external Actions have no guarantee the code that ran yesterday is the code that ran today.
  - The only available protection is to limit what Actions that are allowed to be used in a project to Actions with a good security hygiene.
- Go mod
  - Reference: <https://go.dev/ref/mod#authenticating>
- Dependencies in Go can (default, but alternatives exist) be pulled directly from a git repository. While you may reference the version you want to use by the commit hash,   the Go module system will actually compute a SHA256 hash over the source and use that in the `go.sum` file instead. It is this hash that's used for integrity purposes
  - Thus, even if a malicious commit is created for an existing commit hash, the Go module system would detect this through the SHA256 hash of the source code.
- Deno
  - Dependencies in Deno can (^alternatives exist) be pulled directly from a git repository. As far as I can tell it doesn't rely on the commit hash for integrity ([ref](https://docs.deno.com/runtime/manual/basics/modules/integrity_checking)) and it seems to be using SHA256 ([ref](https://github.com/denoland/deno_lockfile/blob/75e3b2800e3fd6f3d62478a6cf15fd030d91a363/src/lib.rs#L29)), though I'm unsure what it hashes.

## Outcome

- Deno and Go are unaffected. Only the GitHub Actions ecosystem is affected.
- Build a system for GitHub Actions that stores and compares better hashes of the Actions source code (i.e. repository). Hashes could be stored in the repository and checked as a first step of any job.
  - <https://github.com/ericcornelissen/ghasum>
