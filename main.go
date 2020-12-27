package main

import (
	"fmt"
	"github.com/dashtangui/topn/models" 
	"github.com/dashtangui/topn/services" 
	"github.com/dashtangui/topn/clients" 
	"errors"
	"flag"
	"os"
)

const maxPageSize = 100

func main(){
	
	authToken := os.Getenv("GITHUB_AUTH_TOKEN")

	customizeFlagUsage()
	topNFlags := parseCommandFlags()
	
	if _,err:=validateCommandFlags(topNFlags);err!=nil{
		fmt.Printf(err.Error())
		return
	}

	repoRanker := services.TopNRepoRanker{
		TopN: topNFlags.TopN, 
		RepoClient: clients.GitClient,
		PageSize: maxPageSize,
		AuthorizationToken: authToken}

	switch topNFlags.RankedBy {
		case models.RankByOptions.Stars:
			topNByStars:= repoRanker.FetchTopNRepositories(topNFlags.Organization, models.RankByOptions.Stars, topNFlags.Order)
			for _,r:= range topNByStars{
				fmt.Printf("Repo: %v, Count: %v \n",*r.Name, *r.StargazersCount )
			}
		case models.RankByOptions.Forks:
			topNByForks := repoRanker.FetchTopNRepositories(topNFlags.Organization, models.RankByOptions.Forks, topNFlags.Order)
			for _,r := range topNByForks{
				fmt.Printf("Repo: %v, Count: %v \n",*r.Name, *r.ForksCount )
			}
		case models.RankByOptions.PullRequests:
			topNByPullRequests := repoRanker.FetchTopNRepositoriesByPullRequests(topNFlags.Organization, topNFlags.Order)
			for _,r := range topNByPullRequests{
				fmt.Printf("Repo Name: %v, Count: %v \n",r.Name, r.Count )
			}
		case models.RankByOptions.ContributionPercentage:
			topNByContributionPercentage := repoRanker.FetchTopNByContributionPercentage(topNFlags.Organization, topNFlags.Order)
			for _,r:= range topNByContributionPercentage{
				fmt.Printf("Repo Name: %v, Count: %.2f  \n",r.Name, r.Count )
			}
	}	
}

func validateCommandFlags(flags models.TopNFlags)(bool,error){
	if flags.Organization == ""{
		return false, errors.New("Organization is a required flag. Please use -org to provide a GitHub Organization.")
	}

	return true,nil
}

func parseCommandFlags()models.TopNFlags{
	organization := flag.String("org","","Get TopN repos for the provided GitHub org.")
	
	by := flag.String("by",models.RankByOptions.Stars,"Get TopN repos by number of stars,forks,prs,cp.")
	order := flag.String("order",models.OrderByOptions.Desc,"Get TopN repos by in the specified order. Use asc or desc")
	
	topN:= flag.Int("n",1,"TopN will include the top n repositories.")
	pageSize := flag.Int("pageSize",maxPageSize, "Set the perPage to be used when communicating with GitHub API")
	
	flag.Parse()

	topNFlags := models.TopNFlags{
		RankedBy: *by, 
		Organization: *organization, 
		TopN: *topN, 
		PageSize: *pageSize,
		Order: *order}

	return topNFlags
}

func customizeFlagUsage(){
	flag.Usage = func() {
		fmt.Println()
		fmt.Printf("TopN is a command line tool to list the top N repos of a GitHub Organization by number of Stars, Forks, Pull Requests, and Contribution Percentage. \n\n")
		fmt.Println("Set GITHUB_AUTH_TOKEN for better Rate Limits.")
		fmt.Printf("Usage: \n\n")
		fmt.Println("topn -org {org} -n {n} -by {by}")
		flag.PrintDefaults()
	}
}