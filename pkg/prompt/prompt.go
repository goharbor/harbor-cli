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
package prompt

import (
	"errors"
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/context/switch"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	immview "github.com/goharbor/harbor-cli/pkg/views/immutable/select"
	instview "github.com/goharbor/harbor-cli/pkg/views/instance/select"
	lview "github.com/goharbor/harbor-cli/pkg/views/label/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	qview "github.com/goharbor/harbor-cli/pkg/views/quota/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	sview "github.com/goharbor/harbor-cli/pkg/views/scanner/select"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	wview "github.com/goharbor/harbor-cli/pkg/views/webhook/select"
	log "github.com/sirupsen/logrus"
)

func GetRegistryNameFromUser() int64 {
	registryId := make(chan int64)
	go func() {
		response, _ := api.ListRegistries()
		rview.RegistryList(response.Payload, registryId)
	}()

	return <-registryId
}

func GetProjectNameFromUser() (string, error) {
	type result struct {
		name string
		err  error
	}
	resultChan := make(chan result)

	go func() {
		response, err := api.ListAllProjects()
		if err != nil {
			resultChan <- result{"", err}
			return
		}

		if len(response.Payload) == 0 {
			resultChan <- result{"", errors.New("no projects found")}
			return
		}

		name, err := pview.ProjectList(response.Payload)
		if err != nil {
			if err == pview.ErrUserAborted {
				resultChan <- result{"", errors.New("user aborted project selection")}
			} else {
				resultChan <- result{"", fmt.Errorf("error during project selection: %w", err)}
			}
			return
		}

		resultChan <- result{name, nil}
	}()

	res := <-resultChan
	return res.name, res.err
}

func GetRepoNameFromUser(projectName string) string {
	repositoryName := make(chan string)

	go func() {
		response, err := api.ListRepository(projectName, false)
		if err != nil {
			log.Fatal(err)
		}
		repoView.RepositoryList(response.Payload, repositoryName)
	}()

	return <-repositoryName
}

// complete the function
func GetReferenceFromUser(repositoryName string, projectName string) string {
	reference := make(chan string)
	go func() {
		response, _ := api.ListArtifact(projectName, repositoryName)
		aview.ListArtifacts(response.Payload, reference)
	}()
	return <-reference
}

func GetUserIdFromUser() int64 {
	userId := make(chan int64)

	go func() {
		response, _ := api.ListUsers()
		uview.UserList(response.Payload, userId)
	}()

	return <-userId
}

func GetImmutableTagRule(projectName string) int64 {
	immutableid := make(chan int64)
	go func() {
		response, _ := api.ListImmutable(projectName)
		immview.ImmutableList(response.Payload, immutableid)
	}()
	return <-immutableid
}

func GetTagFromUser(repoName, projectName, reference string) string {
	tag := make(chan string)
	go func() {
		response, _ := api.ListTags(projectName, repoName, reference)
		tview.ListTags(response.Payload, tag)
	}()
	return <-tag
}

func GetTagNameFromUser() string {
	repoName := make(chan string)

	go func() {
	}()

	return <-repoName
}

func GetScannerIdFromUser() string {
	scannerId := make(chan string)

	go func() {
		response, _ := api.ListScanners()
		sview.ScannerList(response.Payload, scannerId)
	}()

	return <-scannerId
}

func GetWebhookFromUser(projectName string) (models.WebhookPolicy, error) {
	type result struct {
		webhook models.WebhookPolicy
		err     error
	}

	resultChan := make(chan result)

	go func() {
		res, err := api.ListWebhooks(projectName)
		if err != nil {
			resultChan <- result{models.WebhookPolicy{}, err}
			return
		}

		if len(res.Payload) == 0 {
			resultChan <- result{models.WebhookPolicy{}, errors.New("no webhooks found")}
			return
		}

		webhook, err := wview.WebhookList(res.Payload)
		if err != nil {
			if err == wview.ErrUserAborted {
				resultChan <- result{models.WebhookPolicy{}, errors.New("user aborted webhook selection")}
			} else {
				resultChan <- result{models.WebhookPolicy{}, fmt.Errorf("error during webhook selection: %w", err)}
			}
			return
		}
		resultChan <- result{webhook, nil}
	}()

	res := <-resultChan
	return res.webhook, res.err
}

func GetLabelIdFromUser(opts api.ListFlags) int64 {
	labelId := make(chan int64)
	go func() {
		response, _ := api.ListLabel(opts)
		lview.LabelList(response.Payload, labelId)
	}()

	return <-labelId
}

func GetInstanceFromUser() string {
	instanceName := make(chan string)

	go func() {
		response, _ := api.ListInstance()
		instview.InstanceList(response.Payload, instanceName)
	}()

	return <-instanceName
}

func GetQuotaIDFromUser() int64 {
	QuotaID := make(chan int64)

	go func() {
		response, err := api.ListQuota(*&api.ListQuotaFlags{})
		if err != nil {
			log.Errorf("failed to list quota: %v", err)
		}
		qview.QuotaList(response.Payload, QuotaID)
	}()

	return <-QuotaID
}

func GetActiveContextFromUser() (string, error) {
	config, err := utils.GetCurrentHarborConfig()
	if err != nil {
		return "", err
	}
	var cxlist []api.ContextListView
	for _, cred := range config.Credentials {
		cx := api.ContextListView{Name: cred.Name, Username: cred.Username, Server: cred.ServerAddress}
		cxlist = append(cxlist, cx)
	}

	res, err := list.ContextList(cxlist)
	if err != nil {
		return "", err
	}

	return res, nil
}
