package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestGetContributorsInProject(t *testing.T) {
	basicHTTPClient := http.Client{}
	randomProject := project{}
	randomProject.Link = fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", "PointCloudLibrary", "pcl")

	response := getContributorsInProject(&basicHTTPClient, &randomProject)

	if len(response) == 0 {
		t.Error()
	}

}
