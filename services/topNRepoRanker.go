package services

import (
	"github.com/google/go-github/v33/github"
	"github.com/dashtangui/topn/models"
	"container/heap"
	"strings"
	"log"
)

type RepositoryClient interface{
	//GetAllRepositoriesByOrg(string,int,string)([]*github.Repository,error)
	SearchAllRepositoriesByOrg(string,string,string,int,string)([]*github.Repository,error)
	SearchAllPRsByOrg(string,int,string)([]*github.Issue,error)
}

type TopNRepoRanker struct{
	TopN int
	RepoClient RepositoryClient
	PageSize int
	AuthorizationToken string
}

var sortBy = struct{
	Forks string
	Stars string
}{
	"forks",
	"stars",
}

var orderBy = struct{
	Desc string
	Asc string
}{
	"desc",
	"asc",
}

func(repoRanker TopNRepoRanker) FetchAllRepositories(organization string, by string, order string)([]*github.Repository){
	searchResults,err:= repoRanker.RepoClient.SearchAllRepositoriesByOrg(organization,by,order,repoRanker.PageSize, repoRanker.AuthorizationToken)
	if err != nil{
		//log.Printf("ERROR: %v", err)
		return nil
	}
	return searchResults
}

func (repoRanker TopNRepoRanker) FetchAllPullRequestsInOrg(organization string)[]*github.Issue{
	issues, err := repoRanker.RepoClient.SearchAllPRsByOrg(organization, repoRanker.PageSize, repoRanker.AuthorizationToken)
	if err!=nil{
		log.Printf("ERROR: %v \n", err)
		return nil
	}
	return issues
}

func(repoRanker TopNRepoRanker) FetchTopNRepositories(organization string, by string, order string)([]*github.Repository){
	allRepos := repoRanker.FetchAllRepositories(organization, by, order)
    if len(allRepos) > repoRanker.TopN{
		topNDesc := allRepos[0: repoRanker.TopN]
		return topNDesc
	}
    return allRepos
}

func(repoRanker TopNRepoRanker) FetchTopNRepositoriesByPullRequests(organization string, order string)[]models.RepositoryGroup{
	allPRs := repoRanker.FetchAllPullRequestsInOrg(organization)
	
    prsByRepo := mapReducePRsByRepository(allPRs)

	repoHeap := &RepositoryHeap{}
	heap.Init(repoHeap)

	for repo, count := range prsByRepo{
		if repoHeap.Len() < repoRanker.TopN{
			heap.Push(repoHeap, &RepositoryNode{Name: repo, Ranking: float64(count)})
		}else if float64(count) > (*repoHeap)[0].Ranking{
			heap.Push(repoHeap, &RepositoryNode{Name: repo, Ranking: float64(count)})
			heap.Pop(repoHeap)
		}
	}
	reposInOrder := inOrder(repoHeap, order, "Pull Requests")

	return reposInOrder
}

func(repoRanker TopNRepoRanker) FetchTopNByContributionPercentage(organization string, order string)[]models.RepositoryGroup{
	allReposByOrg:= repoRanker.FetchAllRepositories(organization,sortBy.Forks,order)
	
	forksInRepos := make(map[string]int)
	for _,repo := range allReposByOrg{
		forksInRepos[*repo.Name]=*repo.ForksCount
	}

	allPRsByOrg:= repoRanker.FetchAllPullRequestsInOrg(organization)
	prsByRepo := mapReducePRsByRepository(allPRsByOrg)

	repoHeap := &RepositoryHeap{}
	heap.Init(repoHeap)

	for repo, prsCount := range prsByRepo{
		if forksCount, exist := forksInRepos[repo]; exist{
			contrib := float64(prsCount)/float64(forksCount)
			if repoHeap.Len() < repoRanker.TopN{
				heap.Push(repoHeap, &RepositoryNode{Name: repo, Ranking: contrib})
			}else if contrib > (*repoHeap)[0].Ranking{
				heap.Push(repoHeap, &RepositoryNode{Name: repo, Ranking: contrib})
				heap.Pop(repoHeap)
			}
		}
	}

	reposInOrder := inOrder(repoHeap, order, "Contribution Percentage")

	return reposInOrder
 }

func mapReducePRsByRepository(allPRs []*github.Issue)map[string]int{
	prsByRepo := make(map[string]int) 
	
	for _,pr := range allPRs {
		if pr.RepositoryURL !=nil{
			repoUrlElems := strings.Split(*pr.RepositoryURL,"/")
			issueRepo := repoUrlElems[len(repoUrlElems)-1]
			prsByRepo[issueRepo]++
		}
	}

	return prsByRepo
}

func inOrder(repoHeap *RepositoryHeap, order string, aggregation string)[]models.RepositoryGroup{	
	topN := make([]models.RepositoryGroup, repoHeap.Len())

	if order == models.OrderByOptions.Asc{
		for cursor:= 0; cursor < len(topN); cursor++ {
			minElement := heap.Pop(repoHeap)	
			topN[cursor]= models.RepositoryGroup{
				Name: minElement.(*RepositoryNode).Name,
				Aggregation: aggregation,
			    Count: minElement.(*RepositoryNode).Ranking, }
		}
	}else{
		for cursor:= len(topN)-1; cursor >= 0; cursor-- {
			minElement := heap.Pop(repoHeap)	
			topN[cursor]= models.RepositoryGroup{
				Name: minElement.(*RepositoryNode).Name,
				Aggregation: aggregation,
			    Count: minElement.(*RepositoryNode).Ranking, }
		}
	}
	return topN
}

