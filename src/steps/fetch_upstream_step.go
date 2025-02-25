package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoOpStep
	BranchName string
}

func (step *FetchUpstreamStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	return repo.Logging.FetchUpstream(step.BranchName)
}
