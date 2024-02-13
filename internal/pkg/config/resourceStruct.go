package config

import "github.com/go-openapi/strfmt"

var (
	OutputType string
	WideOutput bool
)

type ListArtifactOptions struct {
	ProjectName string
	RepositoryName string
}

type ListArtifactRespTable struct {
	Digest string
	Tag string
	Size string
	Type string
	PushTime string
	ArtifactID string
	Details []ListArtifactRespOtype
}

type ListArtifactRespOtype struct {
	Labels any
	Id int64
	PullTime string
	RepositoryID int64
	ScanOverview any
}

type CreateProjectOptions struct {
	ProjectName  string
	Public       bool
	RegistryID   int64
	StorageLimit int64
}

type DeleteProjectOptions struct {
	ProjectNameOrID string
}

type ListProjectOptions struct {
	Name       string
	Owner      string
	Page       int64
	PageSize   int64
	Public     bool
	Q          string
	Sort       string
}

type ListProjectRespTable struct {
	ProjectName string
	AccessLevel string
	Repositories int64
	CreationTime strfmt.DateTime
	Details []ListProjectRespOtype
}

type ListProjectRespOtype struct {
	CVEAllowlist any
	ProjectID int32
	OwnerName string
	Metadata any
}

type GetProjectOptions struct {
	ProjectNameOrID string
}


type GetProjectRespTable struct {
	ProjectName string
	AccessLevel string
	Repositories int64
	CreationTime strfmt.DateTime
	Details []GetProjectRespOtype
}

type GetProjectRespOtype struct {
	CVEAllowlist any
	ProjectID int32
	OwnerName string
	Metadata any
}
type CreateRegistrytOptions struct {
	Name        string
	Type       string
	Url         string
	Description string
	Insecure    bool
	Credential  struct {
		AccessKey    string
		AccessSecret string
		Type        string
	}
}

type DeleteRegistryOptions struct {
	Id int64
}

type ListRegistryOptions struct {
	Page     int64
	PageSize int64
	Q        string
	Sort     string
}

type ListRegistryRespTable struct {
	Name string
	Status string
	Type string
	Url string
	CreationTime strfmt.DateTime
}

type UpdateRegistrytOptions struct {
	Id          int64
	Name        string
	Type       string
	Url         string
	Description string
	Insecure    bool
	Credential  struct {
		AccessKey    string
		AccessSecret string
		Type        string
	}
}

type GetRegistryOptions struct {
	Id int64
}

type GetRegistryRespTable struct {
	Name string
	Status string
	Type string
	Url string
	CreationTime strfmt.DateTime
}

type UpdateRepositoryOptions struct {
	ProjectName string
	RepositoryName string
	Repository struct {
		Id int64
		ProjectId int64
		Name string
		Description string
		ArtifactCount int64
		PullCount int64

	}
}

type LoginOptions struct {
	Name          string
	ServerAddress string
	Username      string
	Password      string
}
