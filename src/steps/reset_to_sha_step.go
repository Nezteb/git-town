package steps

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
)

// ResetToShaStep undoes all commits on the current branch
// all the way until the given SHA.
type ResetToShaStep struct {
	NoOpStep
	Hard bool
	Sha  string
}

func (step *ResetToShaStep) Run(repo *git.ProdRepo, driver hosting.Driver) error {
	currentSha, err := repo.Silent.CurrentSha()
	if err != nil {
		return err
	}
	if step.Sha == currentSha {
		return nil
	}
	return repo.Logging.ResetToSha(step.Sha, step.Hard)
}
