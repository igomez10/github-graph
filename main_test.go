package main

import (
	"testing"
)

func TestGetOrganizationContributorsOrganizations(t *testing.T) {
	randomOrganizationName := "PointCloudLibrary"

	response := getOrganizationContributorsOrganizations(randomOrganizationName)
	if len(response.Organization.MembersWithRole.Edges) == 0 {
		t.Error()
	}

}
