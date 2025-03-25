package mbpe

import "container/heap"

type Queue []Merge

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i, j int) bool {
	return q[j].Less(q[i])
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x any) {
	*q = append(*q, x.(Merge))
}

func (q *Queue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[0 : n-1]

	return item
}

func NewQueue(pairs []Merge) *Queue {
	q := make(Queue, len(pairs))

	for i, pair := range pairs {
		q[i] = pair
	}

	heap.Init(&q)

	return &q
}
