package models


var RankByOptions = struct{
	Stars string
	Forks string
	PullRequests string
	ContributionPercentage string
}{
	"stars",
	"forks",
	"prs",
	"cp",
}

var OrderByOptions = struct{
	Desc string
	Asc string
}{
	"desc",
	"asc",
}

type TopNFlags struct{
	Organization string
	RankedBy string
	Order string
	TopN int
	PageSize int
}

type RepositoryGroup struct{
	Name string
	Aggregation string
	Count float64
}
