package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	_ "gopkg.in/cq.v1"
)

func main() {
	flag.Parse()
	neo4jHost := flag.String("host", "http://neo4j:password@localhost:7474", "neo4j host")

	graphRepresentation := traverseOrganization("kubernetes")
	printOrganizationGraph(graphRepresentation)
	//TODO publish nodes to neo4j

	db, err := openNeo4jConnection(*neo4jHost)
	if err != nil {
		fmt.Println(err)
	} else {
		postGraphToNeo4j(db, graphRepresentation)
	}

}

func postGraphToNeo4j(db *sql.DB, graphRepresentation map[string][]*member) {
	//statement := "{\n  \"statements\" : [ {\n    \"statement\" : \"MERGE (n:Company {name:'" + symbolOrigin + "'}) MERGE (test2:Company {name:'" + newcompany.Name + "'}) MERGE (n)-[:isin]->(test2)\"\n\t}\n]\n}"

	for currentOrganization, membersInOrganization := range graphRepresentation {

		for _, member := range membersInOrganization {
			//query := fmt.Sprintf("MERGE (n:Organization { name:'%s' } ) ;", currentOrganization)

			createOrganization := fmt.Sprintf("MERGE (newOrganization:Organization {name:'%s'})", strings.ToLower(currentOrganization))
			createMember := fmt.Sprintf("MERGE (newUser:User {name:'%s'})", member.Login)
			linkOrganizationToUser := fmt.Sprintf("MERGE (newUser)-[:isMember]->(newOrganization);")

			transactionBuilder := strings.Builder{}
			transactionBuilder.WriteString(createOrganization)
			transactionBuilder.WriteString(createMember)
			transactionBuilder.WriteString(linkOrganizationToUser)

			finalTransaction := transactionBuilder.String()

			response, err := db.Exec(finalTransaction)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("SUCCESS ", response)
			}
		}

	}

	db.Close()
}

func openNeo4jConnection(neo4jURL string) (*sql.DB, error) {

	db, err := sql.Open("neo4j-cypher", neo4jURL)
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Connection established with db")
	}

	//open a connection

	return db, nil
}

func traverseOrganization(visitedOrganization string) map[string][]*member {

	kubeRelatedOrganizations := getOrganizationContributorsOrganizations(visitedOrganization)
	organizationsGraph := buildGraph(kubeRelatedOrganizations, visitedOrganization)
	return organizationsGraph
}

func buildGraph(rawAPIResponse OrganizationContributorsOrganizations, initialOrganization string) map[string][]*member {

	//prettyPrint, err := json.MarshalIndent(rawAPIResponse, "", "    ")
	//checkError(err)
	organizationsGraph := make(map[string][]*member)

	/// Add current members to list of members under visited organization
	for _, organizationConnection := range rawAPIResponse.Organization.MembersWithRole.Edges {
		//traverse over members of current organization and add them to initial organization list of members
		currentMember := member(organizationConnection.Node)
		organizationsGraph[initialOrganization] = append(organizationsGraph[initialOrganization], &currentMember)

		// traverse others organizations currentMember belongs to
		for _, currentOrganization := range currentMember.Organizations.Edges {
			// organization already existes, append current member to list of members
			if _, ok := organizationsGraph[currentOrganization.Node.Name]; ok {
				updatedArrayInGraph := append(organizationsGraph[currentOrganization.Node.Name], &currentMember)
				organizationsGraph[currentOrganization.Node.Name] = updatedArrayInGraph
			} else {
				//initialize organization
				organizationsGraph[currentOrganization.Node.Name] = []*member{&currentMember}
			}
		}
	}
	return organizationsGraph
}

func printOrganizationGraph(givenGraph map[string][]*member) {
	for organization, members := range givenGraph {
		fmt.Println(organization)
		for i, currMember := range members {
			fmt.Printf("\t %d. %s \n", i+1, currMember.Login)
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
