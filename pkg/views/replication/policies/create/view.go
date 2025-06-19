// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package create

// CopyByChunk *bool `json:"copy_by_chunk,omitempty"`
// CreationTime strfmt.DateTime `json:"creation_time,omitempty"`
// Deletion bool `json:"deletion,omitempty"`
// Description string `json:"description,omitempty"`
// DestNamespace string `json:"dest_namespace,omitempty"`
// // Specify how many path components will be replaced by the provided destination namespace.
// // The default value is -1 in which case the legacy mode will be applied.
// DestNamespaceReplaceCount *int8 `json:"dest_namespace_replace_count,omitempty"`
// DestRegistry *Registry `json:"dest_registry,omitempty"`
// Enabled bool `json:"enabled,omitempty"`
// // The replication policy filter array.
// Filters []*ReplicationFilter `json:"filters"`
// ID int64 `json:"id,omitempty"`
// Name string `json:"name,omitempty"`
// Override bool `json:"override,omitempty"`
// ReplicateDeletion bool `json:"replicate_deletion,omitempty"`
// Speed *int32 `json:"speed,omitempty"`
// SrcRegistry *Registry `json:"src_registry,omitempty"`
// Trigger *ReplicationTrigger `json:"trigger,omitempty"`
// UpdateTime strfmt.DateTime `json:"update_time,omitempty"`

// struct to hold RPolicy options
type RPolicyCreateView struct {
	Name            string
	Description     string
	ReplicationMode string
	SrcRegistry     string
	DestRegistry    string
	DestNamespace   string
	FlatteningLevel string
	TriggerMode     string
	BandWidth       int64
	Override        bool
	Enabled         bool
	CopyByChunk     bool
}

// func CreateRegistryView(createView *api.CreateRegView) {
// 	registries, _ := api.GetRegistryProviders()

// 	// Initialize a slice to hold registry options
// 	var registryOptions []RegistryOption

// 	// Iterate over registries to populate registryOptions
// 	for i, registry := range registries {
// 		registryOptions = append(registryOptions, RegistryOption{
// 			ID:   strconv.FormatInt(int64(i), 10),
// 			Name: registry,
// 		})
// 	}

// 	// Initialize a slice to hold select options
// 	var registrySelectOptions []huh.Option[string]

// 	// Iterate over registryOptions to populate registrySelectOptions
// 	for _, option := range registryOptions {
// 		registrySelectOptions = append(
// 			registrySelectOptions,
// 			huh.NewOption(option.Name, option.Name),
// 		)
// 	}

// 	theme := huh.ThemeCharm()
// 	err := huh.NewForm(
// 		huh.NewGroup(
// 			huh.NewSelect[string]().
// 				Title("Select a Registry Provider").
// 				Value(&createView.Type).
// 				Options(registrySelectOptions...).
// 				Validate(func(str string) error {
// 					if str == "" {
// 						return errors.New("registry provider cannot be empty")
// 					}
// 					return nil
// 				}),

// 			huh.NewInput().
// 				Title("Name").
// 				Value(&createView.Name).
// 				Validate(func(str string) error {
// 					if strings.TrimSpace(str) == "" {
// 						return errors.New("name cannot be empty or only spaces")
// 					}
// 					if isVaild := utils.ValidateRegistryName(str); !isVaild {
// 						return errors.New("please enter the correct name format")
// 					}
// 					return nil
// 				}),
// 			huh.NewInput().
// 				Title("Description").
// 				Value(&createView.Description),
// 			huh.NewInput().
// 				Title("URL").
// 				Value(&createView.URL).
// 				Validate(func(str string) error {
// 					if strings.TrimSpace(str) == "" {
// 						return errors.New("url cannot be empty or only spaces")
// 					}
// 					formattedUrl := utils.FormatUrl(str)
// 					if _, err := url.ParseRequestURI(formattedUrl); err != nil {
// 						return errors.New("please enter the correct url format")
// 					}
// 					return nil
// 				}),
// 			huh.NewInput().
// 				Title("Access Key").
// 				Value(&createView.Credential.AccessKey),
// 			huh.NewInput().
// 				Title("Access Secret").
// 				Value(&createView.Credential.AccessSecret),
// 			huh.NewConfirm().
// 				Title("Verify Cert").
// 				Value(&createView.Insecure).
// 				Affirmative("yes").
// 				Negative("no"),
// 		),
// 	).WithTheme(theme).Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
