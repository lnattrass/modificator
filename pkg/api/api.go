package api

import (
	"context"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/google/go-github/v60/github"
	"github.com/lnattrass/modificator/pkg/ghapi"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

func Commit(ctx context.Context, token, owner, repository, branch, path, message string, fileContents []byte, createPR bool, mergePR bool) error {
	api := github.NewClient(nil).WithAuthToken(token)

	defaultBranch, err := ghapi.GetDefaultBranch(ctx, api, owner, repository)
	if err != nil {
		return errors.Wrap(err, "failed to get base branch")
	}

	if _, err := ghapi.CreateOrGetBranch(ctx, api, owner, repository, *defaultBranch.Commit.SHA, branch); err != nil {
		return errors.Wrap(err, "failed to create/ensure branch exists")
	}

	if err := ghapi.PutFile(ctx, api, owner, repository, branch, path, message, fileContents); err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	if !createPR {
		return nil
	}

	pr, err := ghapi.CreatePR(ctx, api, owner, repository, *defaultBranch.Name, branch, message)
	if err != nil {
		return errors.Wrap(err, "failed to create the PR")
	}

	if !mergePR {
		return nil
	}

	if err := ghapi.Merge(ctx, api, owner, repository, *pr.Number, message); err != nil {
		return err
	}

	return nil
}

func Patch(ctx context.Context, token, owner, repository, branch, path, message string, patch []byte, createPR bool, mergePR bool) error {
	api := github.NewClient(nil).WithAuthToken(token)

	// Do this now, to not waste quota
	patchJson, err := yaml.YAMLToJSON(patch)
	if err != nil {
		return errors.Wrap(err, "failed to convert patch to JSON")
	}

	defaultBranch, err := ghapi.GetDefaultBranch(ctx, api, owner, repository)
	if err != nil {
		return errors.Wrap(err, "failed to get base branch")
	}

	if _, err := ghapi.CreateOrGetBranch(ctx, api, owner, repository, *defaultBranch.Commit.SHA, branch); err != nil {
		return errors.Wrap(err, "failed to create/ensure branch exists")
	}

	f, err := ghapi.GetFile(ctx, api, owner, repository, branch, path)
	if err != nil {
		return errors.Wrap(err, "failed to get file")
	}

	content, err := f.GetContent()
	if err != nil {
		return errors.Wrap(err, "failed to get content of file")
	}

	// Apply the patch here:
	sourceJson, err := yaml.YAMLToJSON([]byte(content))
	if err != nil {
		return errors.Wrap(err, "failed to conver to JSON")
	}

	patchedJson, err := jsonpatch.MergePatch(sourceJson, patchJson)
	if err != nil {
		return errors.Wrap(err, "failed to mergepatch")
	}

	patchedYaml, err := yaml.JSONToYAML(patchedJson)
	if err != nil {
		return errors.Wrap(err, "failed to marshal back to YAML")
	}

	// End patch

	if err := ghapi.PutFile(ctx, api, owner, repository, branch, path, message, patchedYaml); err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	if !createPR {
		return nil
	}

	pr, err := ghapi.CreatePR(ctx, api, owner, repository, *defaultBranch.Name, branch, message)
	if err != nil {
		return errors.Wrap(err, "failed to create the PR")
	}

	if !mergePR {
		return nil
	}

	if err := ghapi.Merge(ctx, api, owner, repository, *pr.Number, message); err != nil {
		return err
	}

	return nil
}
