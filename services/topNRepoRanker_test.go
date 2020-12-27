package services

import(
	"testing"
	"github.com/google/go-github/v33/github"
	"github.com/dashtangui/topn/models"
	"reflect"
	)

type gitClientMock struct{}

var GitClientMock gitClientMock

var allReposByStars = []*github.Repository{
	{Name: String("Moon"),ForksCount: Int(1), StargazersCount: Int(10)},
	{Name: String("Urano"),ForksCount: Int(0), StargazersCount: Int(5)}, 
	{Name: String("Jupiter"),ForksCount: Int(25), StargazersCount: Int(3)},
	{Name: String("Saturn"),ForksCount: Int(3), StargazersCount: Int(2)}, 
	{Name: String("Sun"),ForksCount: Int(10), StargazersCount: Int(1)}}


var allReposByFork = []*github.Repository{
	{Name: String("Jupiter"),ForksCount: Int(25), StargazersCount: Int(3)},
	{Name: String("Sun"),ForksCount: Int(10), StargazersCount: Int(1)},
	{Name: String("Saturn"),ForksCount: Int(3), StargazersCount: Int(2)}, 
	{Name: String("Moon"),ForksCount: Int(1), StargazersCount: Int(10)},
	{Name: String("Urano"),ForksCount: Int(1), StargazersCount: Int(5)}}

var allPRsInOrg = []*github.Issue{
	{ID: Int64(1), RepositoryURL: String("someurl/Moon")},
	{ID: Int64(2), RepositoryURL: String("someurl/Moon")},
	{ID: Int64(3), RepositoryURL: String("someurl/Moon")},
	{ID: Int64(4), RepositoryURL: String("someurl/Moon")},
	{ID: Int64(5), RepositoryURL: String("someurl/Urano")},
	{ID: Int64(6), RepositoryURL: String("someurl/Urano")},
	{ID: Int64(7), RepositoryURL: String("someurl/Urano")},
	{ID: Int64(8), RepositoryURL: String("someurl/Sun")},
	{ID: Int64(8), RepositoryURL: String("someurl/Sun")},
	{ID: Int64(10), RepositoryURL: String("someurl/Mars")}}

var allPRsByRepo = []models.RepositoryGroup{
	{Name:"Moon",Aggregation:"Pull Requests",Count:4},
	{Name:"Urano",Aggregation:"Pull Requests",Count:3},
	{Name:"Sun",Aggregation:"Pull Requests",Count:2},
	{Name:"Mars",Aggregation:"Pull Requests",Count:1},
}

func (gitClientMock)SearchAllRepositoriesByOrg(organization string, sort string, order string, perPage int, authToken string)([]*github.Repository,error){
	if sort == "stars"{
		return allReposByStars,nil
	}else{
		return allReposByFork,nil
	}
}

func (gitClientMock)SearchAllPRsByOrg(organization string, perPage int, authToken string)([]*github.Issue,error){
	return allPRsInOrg,nil
}

func TestFetchAllRepositories(t *testing.T){
	repoRanker := TopNRepoRanker{
		TopN: 1,
		RepoClient: GitClientMock}

	result:=repoRanker.FetchAllRepositories("organization","stars","")

	if !reflect.DeepEqual(result, allReposByStars) {
		t.Errorf("FetchAllRepositories returned %+v, want %+v", result, allReposByStars)
	}
}

func TestFetchTopNRepositories(t *testing.T){

	testCases := []struct{
		description string
		topN int
		sort string
		want []*github.Repository
	}{
		{"test_Where_N_Is_1",1,"stars",allReposByStars[:1]},
		{"test_Where_N_Is_Less_Than_Total_Elements",3,"stars",allReposByStars[:3]},
		{"test_Where_N_Is_Equal_To_The_Total_Elements",len(allReposByStars),"stars",allReposByStars},
		{"test_Where_N_Is_Greater_Than_The_Total_Elements",1000,"stars",allReposByStars},
		{"test_Expore_Sort_By_Forks",1000,"forks",allReposByFork},
	}
	
	repoRanker := TopNRepoRanker{}
	repoRanker.RepoClient = GitClientMock

	for _,tc:= range testCases{
		repoRanker.TopN = tc.topN
		result := repoRanker.FetchTopNRepositories("org", tc.sort, "")
		if !reflect.DeepEqual(result, tc.want) {
			t.Errorf("%s. returned %+v, want %+v",tc.description, result, tc.want)
		}
	}
}

func TestFetchTopNRepositoriesByPullRequests(t *testing.T){
	testCases := []struct{
		description string
		topN int
		want []models.RepositoryGroup
	}{
		{"test_Where_N_Is_1", 1, allPRsByRepo[:1]},
		{"test_Where_N_Is_Less_Than_Total_Elements", 3, allPRsByRepo[:3]},
		{"test_Where_N_Is_Equal_Than_Total_Elements", len(allPRsByRepo), allPRsByRepo},
		{"test_Where_N_Is_Greater_Than_Total_Elements", 1000, allPRsByRepo},
	}
	
	repoRanker := TopNRepoRanker{}
	repoRanker.RepoClient = GitClientMock

	for _,tc:= range testCases{
		repoRanker.TopN = tc.topN
		result := repoRanker.FetchTopNRepositoriesByPullRequests("org","")
		if !reflect.DeepEqual(result, tc.want) {
			t.Errorf("%s. returned %+v, want %+v",tc.description, result, tc.want)
		}
	}
}

func TestFetchTopNByContributionPercentage(t *testing.T){
	testCases := []struct{
		description string
		topN int
		want []models.RepositoryGroup
	}{
		{"test_Where_N_Is_Greater_Than_Element", 10, []models.RepositoryGroup{
			{Name: "Moon", Aggregation:"Contribution Percentage", Count: 4}, 
			{Name:"Urano", Aggregation:"Contribution Percentage", Count: 3},
			{Name:"Sun", Aggregation:"Contribution Percentage", Count: 0.2}}},
		{"test_Where_N_Is_Equal_Than_Element", 3, []models.RepositoryGroup{
				{Name: "Moon", Aggregation:"Contribution Percentage", Count: 4}, 
				{Name:"Urano", Aggregation:"Contribution Percentage", Count: 3},
				{Name:"Sun", Aggregation:"Contribution Percentage", Count: 0.2}}},
		{"test_Where_N_Is_1", 1, []models.RepositoryGroup{
			{Name: "Moon", Aggregation:"Contribution Percentage", Count: 4}}},
	}
	repoRanker := TopNRepoRanker{}
	repoRanker.RepoClient = GitClientMock

	for _,tc:= range testCases{
		repoRanker.TopN = tc.topN
		result := repoRanker.FetchTopNByContributionPercentage("org","")
		if !reflect.DeepEqual(result, tc.want) {
			t.Errorf("%s. returned %v, want %v",tc.description, result, tc.want)
		}
	}
}

func String(v string)*string{
	return &v
}
func Int64(v int64)*int64{
	return &v
}
func Int(v int)*int{
	return &v
}