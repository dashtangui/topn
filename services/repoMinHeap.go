package services


type RepositoryNode struct{
	Name string
	Ranking float64
}

type RepositoryHeap []*RepositoryNode

func(heap RepositoryHeap) Len() int{
	return len(heap)
}

func(heap RepositoryHeap) Less(i,j int)bool{
	return heap[i].Ranking < heap[j].Ranking
}

func(heap RepositoryHeap) Swap(i,j int){
	heap[i],heap[j] = heap[j],heap[i]
}

func(heap *RepositoryHeap) Push(x interface{}){
	item := x.(*RepositoryNode)
	*heap = append(*heap, item)
}

func(heap *RepositoryHeap) Pop() interface{}{
	old := *heap
	n:= len(old)
	x:=old[n-1]
	old[n-1]=nil
	*heap = old[0 : n-1]
	return x
}