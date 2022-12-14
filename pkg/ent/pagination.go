// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/suyuan32/simple-admin-file/pkg/ent/file"
)

const errInvalidPage = "INVALID_PAGE"

const (
	listField     = "list"
	pageNumField  = "pageNum"
	pageSizeField = "pageSize"
)

type PageDetails struct {
	Page  uint64 `json:"page"`
	Limit uint64 `json:"limit"`
	Total uint64 `json:"total"`
}

// OrderDirection defines the directions in which to order a list of items.
type OrderDirection string

// Cursor of an edge type.
type Cursor struct {
	ID    uint64
	Value Value
}

const (
	// OrderDirectionAsc specifies an ascending order.
	OrderDirectionAsc OrderDirection = "ASC"
	// OrderDirectionDesc specifies a descending order.
	OrderDirectionDesc OrderDirection = "DESC"
)

// Validate the order direction value.
func (o OrderDirection) Validate() error {
	if o != OrderDirectionAsc && o != OrderDirectionDesc {
		return fmt.Errorf("%s is not a valid OrderDirection", o)
	}
	return nil
}

// String implements fmt.Stringer interface.
func (o OrderDirection) String() string {
	return string(o)
}

func (o OrderDirection) reverse() OrderDirection {
	if o == OrderDirectionDesc {
		return OrderDirectionAsc
	}
	return OrderDirectionDesc
}

func (o OrderDirection) orderFunc(field string) OrderFunc {
	if o == OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

type filePager struct {
	order  *FileOrder
	filter func(*FileQuery) (*FileQuery, error)
}

// FilePaginateOption enables pagination customization.
type FilePaginateOption func(*filePager) error

// FileOrder defines the ordering of File.
type FileOrder struct {
	Direction OrderDirection  `json:"direction"`
	Field     *FileOrderField `json:"field"`
}

// FileOrderField defines the ordering field of File.
type FileOrderField struct {
	field    string
	toCursor func(*File) Cursor
}

// DefaultFileOrder is the default ordering of File.
var DefaultFileOrder = &FileOrder{
	Direction: OrderDirectionAsc,
	Field: &FileOrderField{
		field: file.FieldID,
		toCursor: func(f *File) Cursor {
			return Cursor{ID: f.ID}
		},
	},
}

func newFilePager(opts []FilePaginateOption) (*filePager, error) {
	pager := &filePager{}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultFileOrder
	}
	return pager, nil
}

func (p *filePager) applyFilter(query *FileQuery) (*FileQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

// FilePageList is File PageList result.
type FilePageList struct {
	List        []*File      `json:"list"`
	PageDetails *PageDetails `json:"pageDetails"`
}

func (f *FileQuery) Page(
	ctx context.Context, pageNum uint64, pageSize uint64, opts ...FilePaginateOption,
) (*FilePageList, error) {

	pager, err := newFilePager(opts)
	if err != nil {
		return nil, err
	}

	if f, err = pager.applyFilter(f); err != nil {
		return nil, err
	}

	ret := &FilePageList{}

	ret.PageDetails = &PageDetails{
		Page:  pageNum,
		Limit: pageSize,
	}

	count, err := f.Clone().Count(ctx)

	if err != nil {
		return nil, err
	}

	ret.PageDetails.Total = uint64(count)

	direction := pager.order.Direction
	f = f.Order(direction.orderFunc(pager.order.Field.field))
	if pager.order.Field != DefaultFileOrder.Field {
		f = f.Order(direction.orderFunc(DefaultFileOrder.Field.field))
	}

	f = f.Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize))
	list, err := f.All(ctx)
	if err != nil {
		return nil, err
	}
	ret.List = list

	return ret, nil
}
