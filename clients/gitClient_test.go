package clients

import(
	"testing"
	"github.com/google/go-github/v33/github"
	"time"
	"context"
	"reflect"
	"net/http"
)


func TestSearchAllRepositoriesByOrg(t *testing.T) {
	type args struct {
		organization string
		sort         string
		order        string
		perPage      int
		authToken    string
	}
	tests := []struct {
		name    string
		args    args
		responseStatusCode int
		responseNextPage int
		responseError error
		want    []*github.Repository
		wantErr bool
	}{
		{"SearchAllRepositories_Positive_SinglePage", args{"org","stars","",1,""}, 200, 0, nil, []*github.Repository{{ID: Int64(1)}, {ID: Int64(2)}} ,false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SearchRepositories = func(client *github.Client, ctx context.Context, query string, opts *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {

				searchResults := &github.RepositoriesSearchResult{
					Total:             Int(4),
					IncompleteResults: Bool(false),
					Repositories:      tt.want,
				}
				
				rateLimit := github.Timestamp{}
				rateLimit.Time = time.Now().Add(1*time.Second)
				
				response:= Response(&http.Response{StatusCode: tt.responseStatusCode})
				response.Rate = github.Rate{ Reset: rateLimit}
				response.NextPage = tt.responseNextPage

				return searchResults, response, tt.responseError
			}

			got, err := GitClient.SearchAllRepositoriesByOrg(tt.args.organization, tt.args.sort, tt.args.order, tt.args.perPage, tt.args.authToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchAllRepositoriesByOrg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchAllRepositoriesByOrg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Response(response *http.Response) *github.Response {
	return &github.Response{Response: response}
}

func Bool(v bool)*bool{
	return &v
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