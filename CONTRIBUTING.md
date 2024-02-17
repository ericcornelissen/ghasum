<!-- SPDX-License-Identifier: CC0-1.0 -->

# Contributing Guidelines

The maintainers of `ghasum` welcome contributions and corrections. This includes
improvements to the documentation or code base, tests, bug fixes, and
implementations of new features. We recommend you open an issue before making
any substantial changes so you can be sure your work won't be rejected. But for
small changes, such as fixing a typo, you can open a Pull Request directly.

If you decide to make a contribution, please use the following workflow:

- Fork the repository.
- Create a new branch from the latest `main`.
- Make your changes on the new branch.
- Commit to the new branch and push the commit(s).
- Open a Pull Request against `main`.

## Prerequisites

To be able to contribute you need the following tooling:

- [git];
- [Go] v1.21.5 or later;
- (Recommended) a code editor with [EditorConfig] support;

Or a [OCI] compatible container engine, in which case you can run an ephemeral
development container using the command `go run tasks.go dev-env`  (if you don't
have [Go] installed, manually run the commands from the `TaskDevEnv` function in
the `tasks.go` file).

## Tasks

This project uses a custom Go-based task runner to run common tasks. To get
started run:

```shell
go run tasks.go
```

We recommend configuring the following command alias:

```shell
alias gask='go run tasks.go'
```

[editorconfig]: https://editorconfig.org/
[git]: https://git-scm.com/
[go]: https://go.dev/
[oci]: https://opencontainers.org/
