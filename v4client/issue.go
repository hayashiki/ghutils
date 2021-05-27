package v4client

import (
	"github.com/hayashiki/ghutils/v4client/def"
)

type IssueService struct {
	client *v4Client
}

func (s *IssueService) Get(owner, repo, projectName string) (Response, error){
	const query = def.IssueFragments + def.ProjectCardConnectionFragments + def.ProjectConnectionFragments + `
	query ($owner: String!, $name: String!, $searchString: String!) {
		repository(owner: $owner, name: $name) {
			name
			projects(first: 1, search: $searchString) {
				...projectConnection
			}
		}
	}
	`

	variables := map[string]interface{}{
		"owner":        owner,
		"name":         repo,
		"searchString": projectName,
	}

	var resp Response
	err := s.client.Do(query, variables, &resp)
	return resp, err
}

type Response struct {
	Repository Repository `json:"repository"`
}
