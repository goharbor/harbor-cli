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
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/context/switch"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/constants"
	aview "github.com/goharbor/harbor-cli/pkg/views/artifact/select"
	tview "github.com/goharbor/harbor-cli/pkg/views/artifact/tags/select"
	immview "github.com/goharbor/harbor-cli/pkg/views/immutable/select"
	instview "github.com/goharbor/harbor-cli/pkg/views/instance/select"
	lview "github.com/goharbor/harbor-cli/pkg/views/label/select"
	mview "github.com/goharbor/harbor-cli/pkg/views/member/select"
	pview "github.com/goharbor/harbor-cli/pkg/views/project/select"
	qview "github.com/goharbor/harbor-cli/pkg/views/quota/select"
	rview "github.com/goharbor/harbor-cli/pkg/views/registry/select"
	rexecutions "github.com/goharbor/harbor-cli/pkg/views/replication/execution/select"
	rpolicies "github.com/goharbor/harbor-cli/pkg/views/replication/policies/select"
	rtasks "github.com/goharbor/harbor-cli/pkg/views/replication/task/select"

	repoView "github.com/goharbor/harbor-cli/pkg/views/repository/select"
	robotView "github.com/goharbor/harbor-cli/pkg/views/robot/select"
	sview "github.com/goharbor/harbor-cli/pkg/views/scanner/select"
	uview "github.com/goharbor/harbor-cli/pkg/views/user/select"
	wview "github.com/goharbor/harbor-cli/pkg/views/webhook/select"
	log "github.com/sirupsen/logrus"
)

var (
	listRegistriesFunc       = api.ListRegistries
	listArtifactFunc         = api.ListArtifact
	listImmutableFunc        = api.ListImmutable
	listTagsFunc             = api.ListTags
	listScannersFunc         = api.ListScanners
	listInstanceFunc         = api.ListInstance
	listMembersForPromptFunc = api.ListMembers
	listQuotaFunc            = api.ListQuota
)

func GetRegistryNameFromUser() (int64, error) {
	type result struct {
		id  int64
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listRegistriesFunc()
		if err != nil {
			resultChan <- result{0, err}
			return
		}
		idChan := make(chan int64)
		rview.RegistryList(response.Payload, idChan)
		resultChan <- result{<-idChan, nil}
	}()
	res := <-resultChan
	return res.id, res.err
}

func GetProjectIDFromUser() (int64, error) {
	type result struct {
		id  int64
		err error
	}
	resultChan := make(chan result)

	go func() {
		response, err := api.ListAllProjects()
		if err != nil {
			resultChan <- result{0, err}
			return
		}

		if len(response.Payload) == 0 {
			resultChan <- result{0, errors.New("no projects found")}
			return
		}

		id, err := pview.ProjectListWithId(response.Payload)
		if err != nil {
			if err == pview.ErrUserAborted {
				resultChan <- result{0, errors.New("user aborted project selection")}
			} else {
				resultChan <- result{0, fmt.Errorf("error during project selection: %w", err)}
			}
			return
		}

		resultChan <- result{id, nil}
	}()

	res := <-resultChan

	return res.id, res.err
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

// GetRoleNameFromUser prompts the user to select a role and returns it.
func GetRoleNameFromUser() int64 {
	roleChan := make(chan int64)
	Roles := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
	go func() {
		mview.RoleList(Roles, roleChan)
	}()

	return <-roleChan
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
func GetReferenceFromUser(repositoryName string, projectName string) (string, error) {
	type result struct {
		ref string
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listArtifactFunc(projectName, repositoryName)
		if err != nil {
			resultChan <- result{"", err}
			return
		}
		refChan := make(chan string)
		aview.ListArtifacts(response.Payload, refChan)
		resultChan <- result{<-refChan, nil}
	}()
	res := <-resultChan
	return res.ref, res.err
}

func GetUserIdFromUser() (int64, error) {
	response, err := api.ListUsers()
	if err != nil {
		return 0, err
	}

	id, err := uview.UserList(response.Payload)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetImmutableTagRule(projectName string) (int64, error) {
	type result struct {
		id  int64
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listImmutableFunc(projectName)
		if err != nil {
			resultChan <- result{0, err}
			return
		}
		idChan := make(chan int64)
		immview.ImmutableList(response.Payload, idChan)
		resultChan <- result{<-idChan, nil}
	}()
	res := <-resultChan
	return res.id, res.err
}

func GetTagFromUser(repoName, projectName, reference string) (string, error) {
	type result struct {
		tag string
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listTagsFunc(projectName, repoName, reference)
		if err != nil {
			resultChan <- result{"", err}
			return
		}
		tagChan := make(chan string)
		tview.ListTags(response.Payload, tagChan)
		resultChan <- result{<-tagChan, nil}
	}()
	res := <-resultChan
	return res.tag, res.err
}

func GetScannerIdFromUser() (string, error) {
	type result struct {
		id  string
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listScannersFunc()
		if err != nil {
			resultChan <- result{"", err}
			return
		}
		idChan := make(chan string)
		sview.ScannerList(response.Payload, idChan)
		resultChan <- result{<-idChan, nil}
	}()
	res := <-resultChan
	return res.id, res.err
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

func GetLabelIdFromUser(opts api.ListFlags) (int64, error) {
	type result struct {
		id  int64
		err error
	}
	labelId := make(chan result)
	go func() {
		response, err := api.ListLabel(opts)
		if err != nil {
			labelId <- result{0, err}
			return
		}
		choice, err := lview.LabelList(response.Payload)
		if err != nil {
			labelId <- result{0, err}
			return
		}
		labelId <- result{choice, nil}
	}()

	res := <-labelId
	return res.id, res.err
}

func GetInstanceFromUser() (string, error) {
	type result struct {
		name string
		err  error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listInstanceFunc()
		if err != nil {
			resultChan <- result{"", err}
			return
		}
		nameChan := make(chan string)
		instview.InstanceList(response.Payload, nameChan)
		resultChan <- result{<-nameChan, nil}
	}()
	res := <-resultChan
	return res.name, res.err
}

func GetQuotaIDFromUser() (int64, error) {
	type result struct {
		id  int64
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listQuotaFunc(*&api.ListQuotaFlags{})
		if err != nil {
			resultChan <- result{0, err}
			return
		}
		idChan := make(chan int64)
		qview.QuotaList(response.Payload, idChan)
		resultChan <- result{<-idChan, nil}
	}()
	res := <-resultChan
	return res.id, res.err
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

	res, err := list.ContextList(cxlist, config.CurrentCredentialName)
	if err != nil {
		return "", err
	}

	return res, nil
}

func GetRobotPermissionsFromUser(kind string) ([]models.Permission, error) {
	permissions := make(chan robotView.PermissionSelectResult)
	go func() {
		response, err := api.GetPermissions()
		if err != nil {
			permissions <- robotView.PermissionSelectResult{
				Permissions: nil,
				Err:         err,
			}
			return
		}
		robotView.ListPermissions(response.Payload, kind, permissions)
	}()
	result := <-permissions
	return result.Permissions, result.Err
}

func GetRobotIDFromUser(projectID int64) (int64, error) {
	robotID := make(chan int64)
	var opts api.ListFlags
	if projectID != -1 {
		opts.Q = constants.ProjectQString + strconv.FormatInt(projectID, 10)
	}

	go func() {
		response, err := api.ListRobot(opts)
		if err != nil {
			errorCode := utils.ParseHarborErrorCode(err)
			if errorCode == "403" {
				fmt.Println("Permission denied: (Project) Admin privileges are required to execute this command.")
			} else {
				fmt.Printf("failed to list robots: %v\n", utils.ParseHarborErrorMsg(err))
			}
			close(robotID)
			return
		}
		robotView.ListRobot(response.Payload, robotID)
	}()

	id, ok := <-robotID
	if !ok {
		return 0, errors.New("failed to retrieve robot ID")
	}
	return id, nil
}

func GetReplicationPolicyFromUser() int64 {
	replicationPolicyID := make(chan int64)

	go func() {
		response, err := api.ListReplicationPolicies()
		if err != nil {
			log.Fatal(err)
		}
		rpolicies.ReplicationPoliciesList(response.Payload, replicationPolicyID)
	}()

	return <-replicationPolicyID
}

func GetReplicationExecutionIDFromUser(rpolicyID int64) int64 {
	executionID := make(chan int64)

	go func() {
		response, err := api.ListReplicationExecutions(rpolicyID)
		if err != nil {
			log.Fatal(err)
		}
		if len(response.Payload) == 0 {
			log.Fatal("no replication executions found")
		}
		rexecutions.ReplicationExecutionList(response.Payload, executionID)
	}()

	return <-executionID
}

func GetReplicationTaskIDFromUser(execID int64) int64 {
	executionID := make(chan int64)

	go func() {
		response, err := api.ListReplicationTasks(execID)
		if err != nil {
			log.Fatal(err)
		}
		if len(response.Payload) == 0 {
			log.Fatal("no replication tasks found")
		}
		rtasks.ReplicationTasksList(response.Payload, executionID)
	}()

	return <-executionID
}

// Get GetMemberIDFromUser choosing from list of members
func GetMemberIDFromUser(projectName, memberName string) (int64, error) {
	type result struct {
		id  int64
		err error
	}
	resultChan := make(chan result)
	go func() {
		response, err := listMembersForPromptFunc(projectName, memberName, true)
		if err != nil {
			resultChan <- result{0, err}
			return
		}
		if len(response.Payload) == 0 {
			resultChan <- result{0, nil}
			return
		}
		idChan := make(chan int64)
		mview.MemberList(response.Payload, idChan)
		resultChan <- result{<-idChan, nil}
	}()
	res := <-resultChan
	return res.id, res.err
}

// Get Member Role ID selection from user
func GetRoleIDFromUser() int64 {
	roleID := make(chan int64)
	go func() {
		roles := []string{"Project Admin", "Developer", "Guest", "Maintainer", "Limited Guest"}
		mview.RoleList(roles, roleID)
	}()

	return <-roleID
}
