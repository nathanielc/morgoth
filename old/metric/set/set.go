package set

import (
	metric "github.com/nathanielc/morgoth/metric/types"
)

type Set struct {
	set map[metric.MetricID]struct{}
}

func New(capacity int) *Set {
	set := new(Set)
	set.set = make(map[metric.MetricID]struct{}, capacity)
	return set
}
func (set *Set) Add(item metric.MetricID) {
	set.set[item] = struct{}{}
}

func (set *Set) Has(item metric.MetricID) bool {
	_, has := set.set[item]
	return has
}

func (set *Set) Len() metric.MetricID {
	return metric.MetricID(len(set.set))
}

type EachFunc func(metric.MetricID)

func (set *Set) Each(f EachFunc) {
	for item, _ := range set.set {
		f(item)
	}
}
