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


var SearchRepositories = func(client *github.Client, ctx context.Context, query string, opts *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	return client.Search.Repositories(ctx, query, opts)
}

var SearchIssues = func(client *github.Client, ctx context.Context, query string, opts *github.SearchOptions) (*github.IssuesSearchResult, *github.Response, error) {
	return client.Search.Issues(context.Background(), query, opts)
}


func (gitClient)SearchAllRepositoriesByOrg(organization string, sort string, order string, perPage int, authToken string)([]*github.Repository,error){
	
	opts := &github.SearchOptions{Sort: sort, Order: order, ListOptions: github.ListOptions{PerPage: perPage}}
	query := fmt.Sprintf("org:%s",organization)

	client := GitClient.GetClient(authToken)

	var searchResults []*github.Repository
	for {
		results, resp, err := SearchRepositories(client,context.Background(), query, opts)
		if err != nil{
			if resp.StatusCode == 403{
				now := time.Now()
				willResetIn:=resp.Rate.Reset.Time
				if(!now.After(willResetIn)){
					resetDuration:= willResetIn.Sub(now)
					log.Printf("WARNING: Rate Limit Exceeded. Need to go idle for: %v \n", resetDuration)
					time.Sleep(resetDuration)
					log.Printf("MSG: Waking up now.")
				}
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
		results, resp, err :=SearchIssues(client,context.Background(), query, opts)
		if err != nil{
			if resp.StatusCode == 403{
				now := time.Now()
				willResetIn:=resp.Rate.Reset.Time
				if(!now.After(willResetIn)){
					resetDuration:= willResetIn.Sub(now)
					log.Printf("WARNING: Rate Limit Exceeded. Need to go idle for: %v \n", resetDuration)
					time.Sleep(resetDuration)
					log.Printf("MSG: Waking up now.")
				}
			}else{
				log.Printf("ERROR: %v \n", err)
				return nil, err
			}
		}else{
			searchResults = append(searchResults, results.Issues...)
			if resp.NextPage == 0{
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return searchResults,nil
}