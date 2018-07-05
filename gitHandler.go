package main

import (
	"time"

	"gopkg.in/libgit2/git2go.v27"
)

type priorStateInfo struct {
	branchChanged  bool
	prevBranch     string
	changesStashed bool
}

func branchNameFromPath(dir string) string {
	repo, err := git.OpenRepository(dir)
	if err != nil {
		return ""
	}

	return branchNameFromRepo(repo)
}

func branchNameFromRepo(repo *git.Repository) string {
	ref, err := repo.Head()
	if err != nil {
		return ""
	}

	branch, err := ref.Branch().Name()
	if err != nil {
		return ""
	}

	return branch
}

func prepareRepository(dir, branchName string) *priorStateInfo {
	priorState := &priorStateInfo{
		false, "", false,
	}
	repo, err := git.OpenRepository(dir)
	if err != nil {
		// TODO
	}

	if branchNameFromRepo(repo) == branchName {
		return priorState
	}

	priorState.prevBranch = branchName
	priorState.branchChanged = true
	sig := &git.Signature{
		Name: "Nitely Build Service",
		When: time.Now(),
	}

	// TODO: Check whether there are any changes to make
	_, err = repo.Stashes.Save(sig,
		"Automatic stash to test branch "+branchName, git.StashDefault)
	priorState.changesStashed = true
	if err != nil {
		// TODO
	}

	err = trySwitchToBranch(repo, branchName)
	if err != nil {
		// TODO
	}

	return priorState
}

func restoreRepo(dir string, priorState *priorStateInfo) {
	if !priorState.branchChanged {
		return
	}

	repo, err := git.OpenRepository(dir)
	if err != nil {
		// TODO
	}

	err = trySwitchToBranch(repo, priorState.prevBranch)
	if err != nil {
		// TODO
	}

	if priorState.changesStashed {
		err = repo.Stashes.Pop(0, git.StashApplyOptions{})

		if err != nil {
			// TODO
		}
	}
}

func trySwitchToBranch(repo *git.Repository, branchName string) error {
	// TODO: Fix based on whether it's a remote repo.
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return err
	}

	err = repo.SetHead(branch.Reference.Name())
	if err != nil {
		return err
	}

	opts := &git.CheckoutOpts{}
	err = repo.CheckoutHead(opts)
	return err
}
