package clients

import (
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"log"
	"fmt"
	"context"
	"time"
)

type gitClient struct{}

var GitClient gitClient

func (gitClient)GetClient(authToken string)*github.Client{
	var client *github.Client
	if authToken != ""{
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken:authToken},
		)
		tc := oauth2.NewClient(ctx, ts)
	
		client = github.NewClient(tc)
	}
	client = github.NewClient(nil)
	return client
}

func (gitClient)SearchAllRepositoriesByOrg(organization string, sort string, order string, perPage int, authToken string)([]*github.Repository,error){
	
	opts := &github.SearchOptions{Sort: sort, Order: order, ListOptions: github.ListOptions{PerPage: perPage}}
	query := fmt.Sprintf("org:%s",organization)

	client := GitClient.GetClient(authToken)

	var searchResults []*github.Repository
	for {
		results, resp, err := client.Search.Repositories(context.Background(), query, opts)
		if err != nil{
			if resp.StatusCode == 403{
				now := time.Now()
				resetDuration:= resp.Rate.Reset.Time.Sub(now)
				log.Printf("WARNING: Rate Limit Exceeded. Need to go idle for: %v \n", resetDuration)
				time.Sleep(resetDuration)
				log.Printf("MSG: Waking up now.")
			}else{
				log.Printf("ERROR: %v \n", err)
				return nil, err
			}
		}else{
			searchResults = append(searchResults, results.Repositories...)
			if resp.NextPage == 0{
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return searchResults,nil
}

func (gitClient)SearchAllPRsByOrg(organization string, perPage int, authToken string)([]*github.Issue,error){
	
	opts := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: perPage}}
	
	query := fmt.Sprintf("org:%s type:pr",organization)
		
	client := GitClient.GetClient(authToken)

	var searchResults []*github.Issue
	for {
		results, resp, err := client.Search.Issues(context.Background(), query, opts)
		if err != nil{
			if resp.StatusCode == 403{
				now := time.Now()
				resetDuration:= resp.Rate.Reset.Time.Sub(now)
				log.Printf("WARNING: Rate Limit Exceeded. Need to go idle for: %v \n", resetDuration)
				time.Sleep(resetDuration)
				log.Printf("MSG: Waking up now.")
			}else{
				log.Printf("ERROR: %v \n", err)
				return nil, err
			}
		}else{
			searchResults = append(searchResults, results.Issues...)
			if resp.NextPage == 0{
				log.Println("End Of Pages")
				break
			}
			log.Printf("SearchAll - Page: %v \n", opts.Page)
			opts.Page = resp.NextPage
		}
	}

	return searchResults,nil
}

func (gitClient)GetAllRepositoriesByOrg(organization string, pageSize int, authToken string)([]*github.Repository,error){

	client := GitClient.GetClient(authToken)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: pageSize},
	}
	var allRepos []*github.Repository

	for{
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), organization, nil)
		if err != nil{
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0{
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos,nil
}
