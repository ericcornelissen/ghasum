// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/sumdb/dirhash"
	"gopkg.in/yaml.v3"
)

const (
	exitCodeSuccess = 0
	exitCodeError   = 1
	exitCodeUsage   = 2
	exitCodeFailed  = 3
)

var (
	flagDebug = flag.Bool(
		"debug",
		false,
		"Enable debugging mode",
	)
	flagRepo = flag.String(
		"repo",
		"",
		"The repository to work on",
	)
	flagValidate = flag.Bool(
		"validate",
		false,
		"Validate recorded checksums",
	)
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	if *flagRepo == "" {
		fmt.Println("Provide a GitHub repository name using '-repo'")
		return exitCodeUsage
	}

	// ---------------------------------------------------------------------------

	fmt.Println("[INFO ] Creating temporary working directory")
	wd, err := os.MkdirTemp(os.TempDir(), "ghasum-*")
	if err != nil {
		fmt.Printf("[ERROR] couldn't create temporary directory: %s", err)
		return exitCodeError
	}

	if *flagDebug {
		fmt.Printf("[DEBUG] temporary working directory is: '%s'\n", wd)
	} else {
		defer os.RemoveAll(wd)
	}

	// ---------------------------------------------------------------------------

	actions, err := getGhaDependencies(wd, *flagRepo)
	if err != nil {
		fmt.Printf("[ERROR] couldn't obtain GHA dependencies: %s\n", err)
		return exitCodeError
	}

	// ---------------------------------------------------------------------------

	sums := make([][2]string, len(actions))
	for i, action := range actions {
		sum, err := getRepositoryHash(wd, action)
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			return exitCodeError
		}

		sums[i] = [2]string{action, sum}
	}

	// ---------------------------------------------------------------------------

	if *flagValidate {
		fmt.Println("[INFO ] Validating obtained sums against gha.sum")
		ok, err := validateSums(sums)
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			return exitCodeError
		}

		if !ok {
			return exitCodeFailed
		} else {
			fmt.Println("Validation successful")
		}
	} else {
		fmt.Println("[INFO ] Storing obtained sums in gha.sum")
		if err := storeSums(sums); err != nil {
			fmt.Printf("[ERROR] %s\n", err)
			return exitCodeError
		}
	}

	return exitCodeSuccess
}

// -----------------------------------------------------------------------------

func validateSums(sums [][2]string) (ok bool, err error) {
	rawStored, err := os.ReadFile("gha.sum")
	if err != nil {
		return ok, fmt.Errorf("couldn't read gha.sum file: %v", err)
	}

	lines := strings.Split(string(rawStored), "\n")
	lines = lines[:len(lines)-1]
	if len(lines) != len(sums) {
		return false, nil
	}

	for i, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return ok, fmt.Errorf("invalid gha.sum")
		}
		expectedDep, expectedSum := parts[0], parts[1]

		actualDep, actualSum := sums[i][0], sums[i][1]

		if actualDep != expectedDep || actualSum != expectedSum {
			return false, nil
		}
	}

	return true, nil
}

func storeSums(sums [][2]string) (err error) {
	var sb strings.Builder
	for _, sum := range sums {
		name, hash := sum[0], sum[1]
		sb.WriteString(name)
		sb.WriteRune(' ')
		sb.WriteString(hash)
		sb.WriteRune('\n')
	}

	f, err := os.Create("gha.sum")
	if err != nil {
		return fmt.Errorf("couldn't create gha.sum file: %v", err)
	}

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("couldn't write to gha.sum file: %v", err)
	}

	return nil
}

// -----------------------------------------------------------------------------

func getGhaDependencies(wd string, repoName string) (actions []string, err error) {
	url := repoNameToGitHubUrl(repoName)

	wd = path.Join(wd, ".source")

	fmt.Printf("[INFO ] Cloning source repository (url: '%s')\n", url)
	_, err = git.PlainClone(wd, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return actions, fmt.Errorf("couldn't clone the repository: %v", err)
	}

	wd = path.Join(wd, ".github", "workflows")
	dirEntries, err := os.ReadDir(wd)
	if err != nil {
		return actions, fmt.Errorf("couldn't read workflows directory: %v", err)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		if ext := path.Ext(entry.Name()); ext != ".yml" && ext != ".yaml" {
			continue
		}

		workflowPath := path.Join(wd, entry.Name())

		workflowRaw, err := os.ReadFile(workflowPath)
		if err != nil {
			return actions, fmt.Errorf("couldn't read workflow '%s': %v", workflowPath, err)
		}

		workflow, err := parseGhaWorkflow(workflowRaw)
		if err != nil {
			return actions, fmt.Errorf("couldn't parse workflow '%s': %v", workflowPath, err)
		}

		for _, job := range workflow.Jobs {
			for _, step := range job.Steps {
				if step.Uses != "" {
					if !slices.Contains(actions, step.Uses) {
						actions = append(actions, step.Uses)
					}
				}
			}
		}
	}

	sort.Strings(actions)
	return actions, nil
}

func parseGhaWorkflow(data []byte) (w workflow, err error) {
	if err = yaml.Unmarshal(data, &w); err != nil {
		return w, fmt.Errorf("couldn't parse workflow: %v", err)
	}

	return w, nil
}

type workflow struct {
	Jobs map[string]job `yaml:"jobs"`
}

type job struct {
	Steps []step `yaml:"steps"`
}

type step struct {
	Uses string `yaml:"uses"`
}

// -----------------------------------------------------------------------------

func getRepositoryHash(wd string, action string) (sum string, err error) {
	if strings.HasPrefix(action, ".") {
		return "n/a", nil
	}

	parts := strings.Split(action, "@")
	repo, ref := parts[0], parts[1]
	repo = strings.Join(strings.Split(repo, "/")[0:2], "/")
	url := repoNameToGitHubUrl(repo)
	wd = path.Join(wd, action)

	fmt.Printf("[INFO ] Cloning dependency repository (url: '%s')\n", url)
	r, err := git.PlainClone(wd, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return sum, fmt.Errorf("couldn't clone the repository: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return sum, fmt.Errorf("couldn't obtain worktree: %v", err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(ref),
	})
	if err != nil {
		return sum, fmt.Errorf("couldn't checkout specific ref (%s): %v", ref, err)
	}

	// NOTE: Remove the .git index to ensure the directory hash is reproducible.
	// The contents of this file my be different every clone.
	if err := os.RemoveAll(path.Join(wd, ".git")); err != nil {
		return sum, fmt.Errorf("couldn't remove git index: %v", err)
	}

	// ---------------------------------------------------------------------------

	sum, err = dirhash.HashDir(wd, "", dirhash.DefaultHash)
	if err != nil {
		return sum, fmt.Errorf("couldn't compute hash: %v", err)
	}

	return sum, nil
}

// -----------------------------------------------------------------------------

func repoNameToGitHubUrl(name string) (url string) {
	return fmt.Sprintf("https://github.com/%s", name)
}
