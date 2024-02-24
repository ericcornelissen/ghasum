// Copyright 2024 Eric Cornelissen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ericcornelissen/ghasum/internal/cache"
	"github.com/ericcornelissen/ghasum/internal/ghasum"
)

func cmdVerify(argv []string) error {
	var (
		flags       = flag.NewFlagSet(cmdNameVerify, flag.ContinueOnError)
		flagCache   = flags.String(flagNameCache, "", "")
		flagNoCache = flags.Bool(flagNameNoCache, false, "")
	)

	flags.Usage = func() { fmt.Fprintln(os.Stderr) }
	if err := flags.Parse(argv); err != nil {
		return errUsage
	}

	args := flags.Args()
	if len(args) > 1 {
		return errUsage
	}

	target, err := getTarget(args)
	if err != nil {
		return err
	}

	c, err := cache.New(*flagCache, *flagNoCache)
	if err != nil {
		return errors.Join(errCache, err)
	}

	stat, err := os.Stat(target)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	var workflow string
	if !stat.IsDir() {
		repo := path.Join(path.Dir(target), "..", "..")
		workflow, _ = filepath.Rel(repo, target)
		target = repo
	}

	cfg := ghasum.Config{
		Repo:     os.DirFS(target),
		Path:     target,
		Workflow: workflow,
		Cache:    c,
	}

	problems, err := ghasum.Verify(&cfg)
	if err != nil {
		return errors.Join(errUnexpected, err)
	}

	if cnt := len(problems); cnt > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d problems(s) occurred during validation:\n", cnt))
		for _, problem := range problems {
			sb.WriteString(fmt.Sprintf("  %s\n", problem))
		}

		return errors.Join(errFailure, errors.New(sb.String()))
	}

	fmt.Println("Ok")
	return nil
}

func helpVerify() string {
	return `usage: ghasum verify [flags] [target]

Verify the Actions in the target against the stored checksums. If no target is
provided it will default to the current working directory. If the checksums do
not match this command will error with a non-zero exit code.

The target can be either a directory or a file. If it is a directory it must be
the root of a repository (that is, it should contain the .github directory). In
this case checksums will be verified for every workflow in the repository. If it
is a file it must be a workflow file in a repository. In this case checksums
will be verified only for the given workflow.

If ghasum is not yet initialized this command errors (see "ghasum help init").

The available flags are:

    -cache dir
        The location of the cache directory. This is where ghasum stores and
        looks up repositories it needs.
        Defaults to a directory named .ghasum in the user's home directory.
    -no-cache
        Disable the use of the cache. Makes the -cache flag ineffective.`
}
