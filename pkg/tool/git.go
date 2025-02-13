// Copyright The Helm Authors
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

package tool

import (
	"fmt"
	"strings"

	"github.com/helm/chart-testing/v3/pkg/exec"
)

type Git struct {
	exec exec.ProcessExecutor
}

func NewGit(exec exec.ProcessExecutor) Git {
	return Git{
		exec: exec,
	}
}

func (g Git) FileExistsOnBranch(commit string, file string) bool {
	fileSpec := fmt.Sprintf("%s:%s", commit, file)
	_, err := g.exec.RunProcessAndCaptureOutput("git", "cat-file", "-e", fileSpec)
	return err == nil
}

func (g Git) AddWorktree(path string, ref string) error {
	return g.exec.RunProcess("git", "worktree", "add", path, ref)
}

func (g Git) RemoveWorktree(path string) error {
	return g.exec.RunProcess("git", "worktree", "remove", path)
}

func (g Git) Show(commit string, file string) (string, error) {
	fileSpec := fmt.Sprintf("%s:%s", commit, file)
	return g.exec.RunProcessAndCaptureOutput("git", "show", fileSpec)
}

func (g Git) MergeBase(commit1 string, commit2 string) (string, error) {
	return g.exec.RunProcessAndCaptureOutput("git", "merge-base", commit1, commit2)
}

func (g Git) ListChangedFilesInDirs(commit string, dirs ...string) ([]string, error) {
	changedChartFilesString, err :=
		g.exec.RunProcessAndCaptureOutput("git", "diff", "--find-renames", "--name-only", commit, "--", dirs)
	if err != nil {
		return nil, fmt.Errorf("failed creating diff: %w", err)
	}
	if changedChartFilesString == "" {
		return nil, nil
	}
	return strings.Split(changedChartFilesString, "\n"), nil
}

func (g Git) GetURLForRemote(remote string) (string, error) {
	return g.exec.RunProcessAndCaptureOutput("git", "ls-remote", "--get-url", remote)
}

func (g Git) ValidateRepository() error {
	_, err := g.exec.RunProcessAndCaptureOutput("git", "rev-parse", "--is-inside-work-tree")
	return err
}
