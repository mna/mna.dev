package gitlab

const (
	baseURL = "https://gitlab.com/api/v4/"
)

// curl -H Private-Token:"${GITLAB_API_TOKEN}" https://gitlab.com/api/v4/users/___mna___/projects?visibility=public | jq .
// https://docs.gitlab.com/ee/api/projects.html#list-user-projects
