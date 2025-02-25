package cmd

import (
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/steps"
	"github.com/spf13/cobra"
)

type pruneBranchesConfig struct {
	initialBranch                            string
	localBranchesWithDeletedTrackingBranches []string
	mainBranch                               string
}

func pruneBranchesCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "prune-branches",
		Short: "Deletes local branches whose tracking branch no longer exists",
		Long: `Deletes local branches whose tracking branch no longer exists

Deletes branches whose tracking branch no longer exists from the local repository.
This usually means the branch was shipped or killed on another machine.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := createPruneBranchesConfig(repo)
			if err != nil {
				cli.Exit(err)
			}
			stepList, err := createPruneBranchesStepList(config, repo)
			if err != nil {
				cli.Exit(err)
			}
			runState := runstate.New("prune-branches", stepList)
			err = runstate.Execute(runState, repo, nil)
			if err != nil {
				cli.Exit(err)
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			if err := validateIsConfigured(repo); err != nil {
				return err
			}
			return repo.Config.ValidateIsOnline()
		},
	}
}

func createPruneBranchesConfig(repo *git.ProdRepo) (pruneBranchesConfig, error) {
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return pruneBranchesConfig{}, err
	}
	if hasOrigin {
		err = repo.Logging.Fetch()
		if err != nil {
			return pruneBranchesConfig{}, err
		}
	}
	initialBranch, err := repo.Silent.CurrentBranch()
	if err != nil {
		return pruneBranchesConfig{}, err
	}
	localBranchesWithDeletedTrackingBranches, err := repo.Silent.LocalBranchesWithDeletedTrackingBranches()
	if err != nil {
		return pruneBranchesConfig{}, err
	}
	result := pruneBranchesConfig{
		initialBranch:                            initialBranch,
		localBranchesWithDeletedTrackingBranches: localBranchesWithDeletedTrackingBranches,
		mainBranch:                               repo.Config.MainBranch(),
	}
	return result, nil
}

func createPruneBranchesStepList(config pruneBranchesConfig, repo *git.ProdRepo) (runstate.StepList, error) {
	initialBranch := config.initialBranch
	result := runstate.StepList{}
	for _, branchWithDeletedRemote := range config.localBranchesWithDeletedTrackingBranches {
		if initialBranch == branchWithDeletedRemote {
			result.Append(&steps.CheckoutBranchStep{BranchName: config.mainBranch})
		}
		parent := repo.Config.ParentBranch(branchWithDeletedRemote)
		if parent != "" {
			for _, child := range repo.Config.ChildBranches(branchWithDeletedRemote) {
				result.Append(&steps.SetParentBranchStep{BranchName: child, ParentBranchName: parent})
			}
			result.Append(&steps.DeleteParentBranchStep{BranchName: branchWithDeletedRemote})
		}
		if repo.Config.IsPerennialBranch(branchWithDeletedRemote) {
			result.Append(&steps.RemoveFromPerennialBranchesStep{BranchName: branchWithDeletedRemote})
		}
		result.Append(&steps.DeleteLocalBranchStep{BranchName: branchWithDeletedRemote})
	}
	err := result.Wrap(runstate.WrapOptions{RunInGitRoot: false, StashOpenChanges: false}, repo)
	return result, err
}
