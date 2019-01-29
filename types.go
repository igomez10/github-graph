package main

// OrganizationContributorsOrganizations  Contains data on an organization members and it's subsequent organizations
type OrganizationContributorsOrganizations struct {
	Organization struct {
		MembersWithRole struct {
			Edges []struct {
				Node struct {
					ID            string `json:"id"`
					Login         string `json:"login"`
					Name          string `json:"name"`
					Organizations struct {
						Edges []struct {
							Node struct {
								Name string `json:"name"`
								ID   string `json:"id"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"organizations" graphql:"organizations(first: 100)"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"membersWithRole" graphql:"membersWithRole(first: 100)"`
	} `json:"organization" graphql:"organization(login: $organizationLogin)"`
}

type member struct {
	ID            string `json:"id"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	Organizations struct {
		Edges []struct {
			Node struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"organizations"`
}
