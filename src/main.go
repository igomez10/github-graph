package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var query struct {
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
						} `json:"organizations" graphql:"organizations(first: 3)"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"membersWithRole" graphql:"membersWithRole(first: 3)"`
		} `json:"organization" graphql:"organization(login: \"kubernetes\")"`
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	prettyPrint, err := json.MarshalIndent(query, "", "    ")
	checkError(err)
	fmt.Print(string(prettyPrint))
}

// func old() {
// 	//var visitedProjects []*project
// 	//var visitedContributors []*contributor

// 	APIURL := "https://api.github.com/repos/%s/%s/contributors"

// 	kubernetesURL := fmt.Sprintf(APIURL, "kubernetes", "kubernetes")
// 	kubernetes := project{Link: kubernetesURL}

// 	httpClient := http.Client{}

// 	kuberntesContributions := getContributorsInProject(&httpClient, &kubernetes)

// 	for _, currentContributor := range kuberntesContributions {
// 		fmcurrentContributor.
// 	}

// 	fmt.Printf("%+v", kuberntesContributions)
// }

// curls a projects contributors
func getContributorsInProject(client *http.Client, projectToVisit *project) githubContributorResponse {
	response, err := client.Get(projectToVisit.Link)
	checkError(err)

	defer func() {
		err := response.Body.Close()
		checkError(err)
	}()
	body, err := ioutil.ReadAll(response.Body)
	var currentProjectContributors githubContributorResponse
	checkError(err)

	err = json.Unmarshal(body, &currentProjectContributors)
	checkError(err)

	return currentProjectContributors

}

func checkError(err error) {
	if err != nil {
		log.Println("ERROR", err)
	}
}
