package api

import (
	"fmt"

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
		return fmt.Errorf("Failed to initialize client context")
	}

	_, err = client.Artifact.DeleteArtifact(ctx, &artifact.DeleteArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		switch err.(type) {
		case *artifact.DeleteArtifactForbidden:
			return fmt.Errorf("Forbidden to delete artifact: %s/%s@%s", projectName, repoName, reference)
		case *artifact.DeleteArtifactInternalServerError:
			return fmt.Errorf("Internal server error occurred while deleting artifact: %s/%s@%s", projectName, repoName, reference)
		case *artifact.DeleteArtifactNotFound:
			return fmt.Errorf("Artifact not found: %s/%s@%s", projectName, repoName, reference)
		case *artifact.DeleteArtifactUnauthorized:
			return fmt.Errorf("Unauthorized to delete artifact: %s/%s@%s", projectName, repoName, reference)
		default:
			return fmt.Errorf("Unknown error occurred while deleting artifact: %v", err)
		}
	}

	log.Infof("Artifact deleted successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// InfoArtifact retrieves information about a specific artifact.
func ViewArtifact(projectName, repoName, reference string) (*artifact.GetArtifactOK, error) {
	ctx, client, err := utils.ContextWithClient()
	var response = &artifact.GetArtifactOK{}
	if err != nil {
		return response, fmt.Errorf("Failed to initialize client context")
	}

	response, err = client.Artifact.GetArtifact(ctx, &artifact.GetArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})

	if err != nil {
		log.Errorf("Failed to get artifact info: ")

		switch err.(type) {
		case *artifact.GetArtifactForbidden:
			return response, fmt.Errorf("Forbidden to retrieve artifact: %s/%s@%s", projectName, repoName, reference)
		case *artifact.GetArtifactInternalServerError:
			return response, fmt.Errorf("Internal server error occurred while retrieving artifact: %s/%s@%s", projectName, repoName, reference)
		case *artifact.GetArtifactNotFound:
			return response, fmt.Errorf("Artifact not found: %s/%s@%s", projectName, repoName, reference)
		case *artifact.GetArtifactUnauthorized:
			return response, fmt.Errorf("Unauthorized to retrieve artifact: %s/%s@%s", projectName, repoName, reference)
		default:
			return response, fmt.Errorf("Unknown error occurred while retrieving artifact info: %v", err)
		}
	}

	return response, nil
}

// RunListArtifact lists all artifacts in a repository.
func ListArtifact(projectName, repoName string, opts ...ListFlags) (artifact.ListArtifactsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return artifact.ListArtifactsOK{}, fmt.Errorf("Failed to initialize client context")
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
		switch err.(type) {
		case *artifact.ListArtifactsBadRequest:
			return artifact.ListArtifactsOK{}, fmt.Errorf("Bad request for listing artifacts: %s/%s", projectName, repoName)
		case *artifact.ListArtifactsInternalServerError:
			return artifact.ListArtifactsOK{}, fmt.Errorf("Internal server error occurred while listing artifacts: %s/%s", projectName, repoName)
		case *artifact.ListArtifactsNotFound:
			return artifact.ListArtifactsOK{}, fmt.Errorf("Artifacts not found: %s/%s", projectName, repoName)
		case *artifact.ListArtifactsUnauthorized:
			return artifact.ListArtifactsOK{}, fmt.Errorf("Unauthorized to list artifacts: %s/%s", projectName, repoName)
		default:
			return artifact.ListArtifactsOK{}, fmt.Errorf("Unknown error occurred while listing artifacts: %v", err)
		}
	}

	return *response, nil
}

// StartScanArtifact initiates a scan on a specific artifact.
func StartScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize client context")
	}

	_, err = client.Scan.ScanArtifact(ctx, &scan.ScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to start scan: ")

		switch err.(type) {
		case *scan.ScanArtifactBadRequest:
			return fmt.Errorf("Bad request for starting scan: %s and %s", repoName, projectName)
		case *scan.ScanArtifactForbidden:
			return fmt.Errorf("Scan already in progress for artifact: %s and %s", repoName, projectName)
		case *scan.ScanArtifactInternalServerError:
			return fmt.Errorf("Internal server error occurred while starting scan: %s and %s", repoName, projectName)
		case *scan.ScanArtifactNotFound:
			return fmt.Errorf("Artifacts %s/%s not found", projectName, repoName)
		case *scan.ScanArtifactUnauthorized:
			return fmt.Errorf("Unauthorized to start scan: %s and %s", repoName, projectName)
		default:
			return fmt.Errorf("Unknown error occurred while starting scan: %s and %s: %v", repoName, projectName, err)
		}
	}

	log.Infof("Scan started successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// StopScanArtifact stops a scan on a specific artifact.
func StopScanArtifact(projectName, repoName, reference string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize client context")
	}

	_, err = client.Scan.StopScanArtifact(ctx, &scan.StopScanArtifactParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})
	if err != nil {
		log.Errorf("Failed to stop scan: ")

		switch err.(type) {
		case *scan.StopScanArtifactBadRequest:
			return fmt.Errorf("Bad request for stopping scan: %s/%s@%s", projectName, repoName, reference)
		case *scan.StopScanArtifactForbidden:
			return fmt.Errorf("Scan already stopped for artifact: %s/%s@%s", projectName, repoName, reference)
		case *scan.StopScanArtifactInternalServerError:
			return fmt.Errorf("Internal server error occurred while stopping scan: %s/%s@%s", projectName, repoName, reference)
		case *scan.StopScanArtifactNotFound:
			return fmt.Errorf("Artifact not found: %s/%s@%s", projectName, repoName, reference)
		case *scan.StopScanArtifactUnauthorized:
			return fmt.Errorf("Unauthorized to stop scan: %s/%s@%s", projectName, repoName, reference)
		default:
			return fmt.Errorf("Unknown error occurred while stopping scan: %v", err)
		}
	}

	log.Infof("Scan stopped successfully: %s/%s@%s", projectName, repoName, reference)
	return nil
}

// DeleteTag deletes a tag from a specific artifact.
func DeleteTag(projectName, repoName, reference, tag string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize client context")
	}

	_, err = client.Artifact.DeleteTag(ctx, &artifact.DeleteTagParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
		TagName:        tag,
	})
	if err != nil {
		log.Errorf("Failed to delete tag: ")

		switch err.(type) {
		case *artifact.DeleteTagForbidden:
			return fmt.Errorf("Forbidden to delete tag: %s/%s@%s:%s", projectName, repoName, reference, tag)
		case *artifact.DeleteTagInternalServerError:
			return fmt.Errorf("Internal server error occurred while deleting tag: %s/%s@%s:%s", projectName, repoName, reference, tag)
		case *artifact.DeleteTagNotFound:
			return fmt.Errorf("Tag not found: %s/%s@%s:%s", projectName, repoName, reference, tag)
		case *artifact.DeleteTagUnauthorized:
			return fmt.Errorf("Unauthorized to delete tag: %s/%s@%s:%s", projectName, repoName, reference, tag)
		default:
			return fmt.Errorf("Unknown error occurred while deleting tag: %v", err)
		}
	}

	log.Infof("Tag deleted successfully: %s/%s@%s:%s", projectName, repoName, reference, tag)
	return nil
}

// ListTags lists all tags of a specific artifact.
func ListTags(projectName, repoName, reference string) (*artifact.ListTagsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return &artifact.ListTagsOK{}, fmt.Errorf("Failed to initialize client context")
	}

	resp, err := client.Artifact.ListTags(ctx, &artifact.ListTagsParams{
		ProjectName:    projectName,
		RepositoryName: repoName,
		Reference:      reference,
	})

	if err != nil {

		log.Errorf("Failed to list tags: ")
		switch err.(type) {
		case *artifact.ListTagsBadRequest:
			return &artifact.ListTagsOK{}, fmt.Errorf("Bad request for listing tags: %s/%s", projectName, repoName)
		case *artifact.ListTagsForbidden:
			return &artifact.ListTagsOK{}, fmt.Errorf("Forbidden to list tags: %s/%s", projectName, repoName)
		case *artifact.ListTagsInternalServerError:
			return &artifact.ListTagsOK{}, fmt.Errorf("Internal server error occurred while listing tags: %s/%s", projectName, repoName)
		case *artifact.ListTagsNotFound:
			return &artifact.ListTagsOK{}, fmt.Errorf("Tags not found for artifact: %s/%s@%s", projectName, repoName, reference)
		case *artifact.ListTagsUnauthorized:
			return &artifact.ListTagsOK{}, fmt.Errorf("Unauthorized to list tags: %s/%s", projectName, repoName)
		default:
			return &artifact.ListTagsOK{}, fmt.Errorf("Unknown error occurred while listing tags: %s/%s: %v", projectName, repoName, err)
		}

	}

	return resp, nil
}

// CreateTag creates a tag for a specific artifact.
func CreateTag(projectName, repoName, reference, tagName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("Failed to initialize client context")
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
		log.Errorf("Failed to create tag: ")
		switch err.(type) {
		case *artifact.CreateTagBadRequest:
			return fmt.Errorf("Bad request for creating tag: %s/%s:%s", repoName, reference, tagName)
		case *artifact.CreateTagForbidden:
			return fmt.Errorf("Forbidden to create tag: %s/%s:%s", repoName, reference, tagName)
		case *artifact.CreateTagInternalServerError:
			return fmt.Errorf("Internal server error occurred while creating tag: %s/%s:%s", repoName, reference, tagName)
		case *artifact.CreateTagNotFound:
			return fmt.Errorf("Artifact not found: %s/%s@%s", projectName, repoName, reference)
		case *artifact.CreateTagUnauthorized:
			return fmt.Errorf("Unauthorized to create tag: %s/%s:%s", repoName, reference, tagName)
		default:
			return fmt.Errorf("Unknown error occurred while creating tag: %s/%s:%s: %v", repoName, reference, tagName, err)
		}
	}
	log.Infof("Tag created successfully: %s/%s@%s:%s", projectName, repoName, reference, tagName)
	return nil
}
