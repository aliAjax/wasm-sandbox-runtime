package worker

import "container/heap"

type PriorityQueue []job

func (q PriorityQueue) Len() int { return len(q) }
func (q PriorityQueue) Less(i, j int) bool {
	if q[i].priority == q[j].priority {
		return q[i].created.Before(q[j].created)
	}
	return q[i].priority > q[j].priority
}
func (q PriorityQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *PriorityQueue) Push(x any)   { *q = append(*q, x.(job)) }
func (q *PriorityQueue) Pop() any     { old := *q; n := len(old); v := old[n-1]; *q = old[:n-1]; return v }
func (q *PriorityQueue) Add(j job)    { heap.Push(q, j) }
func (q *PriorityQueue) Take() job    { return heap.Pop(q).(job) }
