package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scan"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// DeleteArtifact handles the deletion of an artifact.
func DeleteArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Artifact.DeleteArtifact(ctx, &artifact.DeleteArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to delete artifact: %v", err)
		return err
	}

	log.Infof("Artifact deleted successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// InfoArtifact retrieves information about a specific artifact.
func ViewArtifact(projectName, repoName, reference string) (*artifact.GetArtifactOK, error) {
	ctx, client, err := utils.ContextWithClient()
	var response = &artifact.GetArtifactOK{}
	if err != nil {
		return response, err
	}

	response, err = client.Artifact.GetArtifact(ctx, &artifact.GetArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})

	if err != nil {
		log.Errorf("Failed to get artifact info: %v", err)
		return response, err
	}

	return response, nil
}

// RunListArtifact lists all artifacts in a repository.
func ListArtifact(projectName, repoName string, opts ...ListFlags) (artifact.ListArtifactsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return artifact.ListArtifactsOK{}, err
	}

	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}
	response, err := client.Artifact.ListArtifacts(ctx, &artifact.ListArtifactsParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Page:           &listFlags.Page,
		PageSize:       &listFlags.PageSize,
		Q:              &listFlags.Q,
		Sort:           &listFlags.Sort,
	})
	if err != nil {
		return artifact.ListArtifactsOK{}, err
	}

	return *response, nil
}

// StartScanArtifact initiates a scan on a specific artifact.
func StartScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scan.ScanArtifact(ctx, &scan.ScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to start scan: %v", err)
		return err
	}

	log.Infof("Scan started successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// StopScanArtifact stops a scan on a specific artifact.
func StopScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scan.StopScanArtifact(ctx, &scan.StopScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to stop scan: %v", err)
		return err
	}

	log.Infof("Scan stopped successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// DeleteTag deletes a tag from a specific artifact.
func DeleteTag(projectName, repoName, reference, tag string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Artifact.DeleteTag(ctx, &artifact.DeleteTagParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
		TagName:        tag,
	})
	if err != nil {
		log.Errorf("Failed to delete tag: %v", err)
		return err
	}

	log.Infof("Tag deleted successfully: %s/%s@%s:%s", projectName, repoName, reference, tag)
	return nil
}

// ListTags lists all tags of a specific artifact.
func ListTags(projectName, repoName, reference string) (*artifact.ListTagsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return &artifact.ListTagsOK{}, err
	}

	resp, err := client.Artifact.ListTags(ctx, &artifact.ListTagsParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})

	if err != nil {
		log.Errorf("Failed to list tags: %v", err)
		return &artifact.ListTagsOK{}, err
	}

	return resp, nil
}

// CreateTag creates a tag for a specific artifact.
func CreateTag(projectName, repoName, reference, tagName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Artifact.CreateTag(ctx, &artifact.CreateTagParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
		Tag: &models.Tag{
			Name: tagName,
		},
	})
	if err != nil {
		log.Errorf("Failed to create tag: %v", err)
		return err
	}
	log.Infof("Tag created successfully: %s/%s@%s:%s", projectName, repoName, reference, tagName)
	return nil
}
