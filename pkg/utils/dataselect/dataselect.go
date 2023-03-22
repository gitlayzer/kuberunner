package dataselect

import (
	"sort"
	"strings"
	"time"
)

type DataSelector struct {
	GenericDataList   []DataCell
	DataSelectorQuery *DataSelectorQuery
}

type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

type DataSelectorQuery struct {
	FilterQuery *FilterQuery
	Pagination  *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

func (d *DataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *DataSelector) Less(i, j int) bool {
	return d.GenericDataList[i].GetCreation().Before(d.GenericDataList[j].GetCreation())
}

func (d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

func (d *DataSelector) Filter() *DataSelector {
	if d.DataSelectorQuery.FilterQuery.Name == "" {
		return d
	}
	var tmp []DataCell
	for _, v := range d.GenericDataList {
		if strings.Contains(v.GetName(), d.DataSelectorQuery.FilterQuery.Name) {
			tmp = append(tmp, v)
		}
	}
	d.GenericDataList = tmp
	return d
}

func (d *DataSelector) Paginate() *DataSelector {
	if d.DataSelectorQuery.Pagination == nil {
		return d
	}
	var tmp []DataCell
	for i, v := range d.GenericDataList {
		if i >= d.DataSelectorQuery.Pagination.Limit*(d.DataSelectorQuery.Pagination.Page-1) && i < d.DataSelectorQuery.Pagination.Limit*d.DataSelectorQuery.Pagination.Page {
			tmp = append(tmp, v)
		}
	}
	d.GenericDataList = tmp
	return d
}
