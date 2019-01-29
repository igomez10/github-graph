package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {

	visitedOrganization := "kubernetes"

	kubeRelatedOrganizations := getOrganizationContributorsOrganizations(visitedOrganization)

	// var contributors []*struct{}
	//kubeRelatedOrganizations.Organization.MembersWithRole.
	// for _, member := kubeRelatedOrganizations.
	prettyPrint, err := json.MarshalIndent(kubeRelatedOrganizations, "", "    ")
	checkError(err)
	fmt.Print(string(prettyPrint))
	var organizationsGraph map[string][]*member
	fmt.Println(organizationsGraph)

	/// Add current members to list of members under visited organization

	for _, organizationConnection := range kubeRelatedOrganizations.Organization.MembersWithRole.Edges {
		currentMember := member(organizationConnection.Node)
		organizationsGraph[visitedOrganization] = append(organizationsGraph[visitedOrganization], &currentMember)

		for _, currentOrganization := range currentMember.Organizations.Edges {
			if ok, organization := organizationGraph[currentOrganization.Node.Name]; ok {

			}
		}

	}

	for _, organizationConnection := range kubeRelatedOrganizations.Organization.MembersWithRole.Edges {
		currentMember := organizationConnection.Node
		if _, exists := organizationsGraph[currentMember.Login]; exists {
			castedMember := member(currentMember)
			organizationsGraph[currentMember.Login] = append(organizationsGraph[currentMember.Login], &castedMember)
		}
	}
}

func getOrganizationContributorsOrganizations(organizationLogin string) OrganizationContributorsOrganizations {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GRAPHQL_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var relatedOrganizations OrganizationContributorsOrganizations

	graphQLVariables := map[string]interface{}{
		"organizationLogin": githubv4.String(organizationLogin),
	}

	err := client.Query(context.Background(), &relatedOrganizations, graphQLVariables)
	if err != nil {
		fmt.Println("Error", err)
	}
	return relatedOrganizations
}

func checkError(err error) {
	if err != nil {
		log.Println("ERROR", err)
	}
}
