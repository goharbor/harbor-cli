package api

type ListFlags struct {
	ProjectID int64
	Scope     string
	Name      string
	Page      int64
	PageSize  int64
	Q         string
	Sort      string
	Public    bool
}

// CreateView for Registry
type CreateRegView struct {
	Name        string
	Type        string
	Description string
	URL         string
	Credential  RegistryCredential
	Insecure    bool
}

// Credential for Registry
type RegistryCredential struct {
	AccessKey    string `json:"access_key,omitempty"`
	Type         string `json:"type,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
}
