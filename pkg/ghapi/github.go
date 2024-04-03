package ghapi

import (
	"context"
	"net/http"

	"github.com/google/go-github/v60/github"
	"github.com/pkg/errors"
)

var (
	ErrDirectory = errors.New("got directory, but expected a file")
)

func CreatePR(ctx context.Context, api *github.Client, owner, repository, baseBranch, headBranch, message string) (*github.PullRequest, error) {
	newPullReq := &github.NewPullRequest{
		Title:               github.String(message),
		Head:                github.String(headBranch),
		Base:                github.String(baseBranch),
		Body:                github.String(message),
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := api.PullRequests.Create(ctx, owner, repository, newPullReq)
	if err != nil {
		return nil, err
	}

	return pr, nil

}

func PutFile(ctx context.Context, api *github.Client, owner, repository, branchName, path, message string, file []byte) error {
	f, err := GetFile(ctx, api, owner, repository, branchName, path)
	if err != nil {
		return errors.Wrap(err, "failed to get existing sha")
	}

	_, _, err = api.Repositories.CreateFile(ctx, owner, repository, path, &github.RepositoryContentFileOptions{
		Message: &message,
		Branch:  &branchName,
		Content: file,
		SHA:     f.SHA,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	return nil
}

func GetFile(ctx context.Context, api *github.Client, owner, repository, branchName, path string) (*github.RepositoryContent, error) {
	f, d, res, err := api.Repositories.GetContents(ctx, owner, repository, path, &github.RepositoryContentGetOptions{
		Ref: "refs/heads/" + branchName,
	})

	if err != nil {
		if res.StatusCode == 404 {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get sha for existing file")
	}

	if d != nil {
		return nil, errors.Wrapf(ErrDirectory, "can't get sha for path %s as it's a directory", path)
	}

	return f, nil
}

func CreateOrGetBranch(ctx context.Context, api *github.Client, owner, repository, baseCommit string, branchName string) (*github.Branch, error) {
	branch, res, err := api.Repositories.GetBranch(ctx, owner, repository, branchName, 2)
	if err == nil {
		// branch exists already
		return branch, nil
	}

	if res.StatusCode != http.StatusNotFound {
		return nil, errors.Wrap(err, "error getting branch")
	}

	ref := &github.Reference{
		Ref: github.String("refs/heads/" + branchName),
		Object: &github.GitObject{
			SHA: github.String(baseCommit),
		},
	}

	if _, _, err := api.Git.CreateRef(ctx, owner, repository, ref); err != nil {
		return nil, errors.Wrap(err, "failed to create remote ref")
	}

	branch, _, err = api.Repositories.GetBranch(ctx, owner, repository, branchName, 2)
	if err != nil {
		return nil, errors.Wrap(err, "branch was created, but we couldnt read after writing")
	}

	return branch, nil
}

func GetDefaultBranch(ctx context.Context, api *github.Client, owner, repository string) (*github.Branch, error) {
	repo, _, err := api.Repositories.Get(ctx, owner, repository)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get repository")
	}

	defaultBranch, _, err := api.Repositories.GetBranch(ctx, owner, repository, *repo.DefaultBranch, 2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get default branch")
	}
	return defaultBranch, nil
}

func Merge(ctx context.Context, api *github.Client, owner, repository string, prNumber int, message string) error {
	_, _, err := api.PullRequests.Merge(ctx, owner, repository, prNumber, message, &github.PullRequestOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to merge PR")
	}

	return nil
}
