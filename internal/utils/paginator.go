package utils

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	defPageAttribute = "page"
)

// PaginatorResources type
type PaginatorResources struct {
	Total         int64  `json:"total"`
	Offset        int    `json:"offset"`
	PagesNum      int    `json:"pages_num"`
	PagesList     []int  `json:"pages_list"`
	PageNo        int    `json:"page_no"`
	PageLink      string `json:"page_link"`
	PageLinkPrev  string `json:"page_prev"`
	PageLinkNext  string `json:"page_next"`
	PageLinkFirst string `json:"page_first"`
	PageLinkLast  string `json:"page_last"`
	PageHasPages  bool   `json:"page_has_pages"`
	PageHasPrev   bool   `json:"page_has_prev"`
	PageHasNext   bool   `json:"page_has_next"`
	PageActive    bool   `json:"page_active"`
}

// IPaginatorBuilder interface
type IPaginatorBuilder interface {
	Build() *PaginatorResources
}

// PaginatorBuilder type
type PaginatorBuilder struct {
	paginator *Paginator
}

// NewPaginatorBuilder constructor
func NewPaginatorBuilder(p *Paginator) IPaginatorBuilder {
	return &PaginatorBuilder{
		paginator: p,
	}
}

// Build method
func (rcv *PaginatorBuilder) Build() *PaginatorResources {
	resources := &PaginatorResources{}

	resources.Total = rcv.paginator.Nums()
	resources.Offset = rcv.paginator.Offset()
	resources.PagesNum = rcv.paginator.PageNums()
	resources.PagesList = rcv.paginator.Pages()
	resources.PageNo = rcv.paginator.Page()
	resources.PageLink = rcv.paginator.PageLink(resources.PageNo)
	resources.PageLinkPrev = rcv.paginator.PageLinkPrev()
	resources.PageLinkNext = rcv.paginator.PageLinkNext()
	resources.PageLinkFirst = rcv.paginator.PageLinkFirst()
	resources.PageLinkLast = rcv.paginator.PageLinkLast()
	resources.PageHasPages = rcv.paginator.HasPages()
	resources.PageHasPrev = rcv.paginator.HasPrev()
	resources.PageHasNext = rcv.paginator.HasNext()
	resources.PageActive = rcv.paginator.IsActive(resources.PageNo)

	return resources
}

// PaginatorOptions type
type PaginatorOptions struct {
	Page  int
	Limit int
	Order string

	Request *http.Request
}

// NewPaginatorOptions constructor
func NewPaginatorOptions(r *http.Request, page, limit, order string) *PaginatorOptions {
	instance := &PaginatorOptions{
		Page: 1, Limit: 20, Order: "asc", Request: r,
	}
	instance.Page, _ = strconv.Atoi(page)
	instance.Limit, _ = strconv.Atoi(limit)

	if strings.ToLower(order) == "desc" {
		instance.Order = "desc"
	}
	return instance
}

// Paginator type
type Paginator struct {
	Request     *http.Request
	PerPageNums int
	MaxPages    int

	nums      int64
	pageRange []int
	pageNums  int
	page      int
}

// NewPaginator constructor
// Instantiates a paginator struct for the current http request.
func NewPaginator(req *http.Request, limit int, total interface{}) *Paginator {
	p := Paginator{}
	p.Request = req

	if limit <= 0 {
		limit = 10
	}
	p.PerPageNums = limit
	p.SetNums(total)

	return &p
}

// PageNums method
// Returns the total number of pages.
func (p *Paginator) PageNums() int {
	if p.pageNums != 0 {
		return p.pageNums
	}
	pageNums := math.Ceil(float64(p.nums) / float64(p.PerPageNums))

	if p.MaxPages > 0 {
		pageNums = math.Min(pageNums, float64(p.MaxPages))
	}
	p.pageNums = int(pageNums)

	return p.pageNums
}

// Nums method
// Returns the total number of items (e.g. from doing SQL count).
func (p *Paginator) Nums() int64 {
	return p.nums
}

// SetNums method
// Sets the total number of items.
func (p *Paginator) SetNums(nums interface{}) {
	p.nums, _ = toInt64(nums)
}

// Page method
// Returns the current page.
func (p *Paginator) Page() int {
	if p.page != 0 {
		return p.page
	}
	if p.Request.Form == nil {
		p.Request.ParseForm()
	}
	p.page, _ = strconv.Atoi(p.Request.Form.Get(defPageAttribute))

	if p.page > p.PageNums() {
		p.page = p.PageNums()
	}
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}

// Pages method
// Returns a list of all pages.
func (p *Paginator) Pages() []int {
	if p.pageRange == nil && p.nums > 0 {
		var pages []int

		page := p.Page()
		pageNums := p.PageNums()

		switch {
		case page >= pageNums-4 && pageNums > 9:
			start := pageNums - 9 + 1
			pages = make([]int, 9)

			for i := range pages {
				pages[i] = start + i
			}
		case page >= 5 && pageNums > 9:
			start := page - 5 + 1
			pages = make([]int, int(math.Min(9, float64(page+4+1))))

			for i := range pages {
				pages[i] = start + i
			}
		default:
			pages = make([]int, int(math.Min(9, float64(pageNums))))

			for i := range pages {
				pages[i] = i + 1
			}
		}
		p.pageRange = pages
	}
	return p.pageRange
}

// PageLink method
// Returns URL for a given page index.
func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.Request.URL.String())
	values := link.Query()

	if page == 1 {
		values.Del(defPageAttribute)
	} else {
		values.Set(defPageAttribute, strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()

	return link.String()
}

// PageLinkPrev method
// Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.Page() - 1)
	}
	return
}

// PageLinkNext method
// Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.Page() + 1)
	}
	return
}

// PageLinkFirst method
// Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

// PageLinkLast method
// Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(p.PageNums())
}

// HasPrev method
// Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return p.Page() > 1
}

// HasNext method
// Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

// IsActive method
// Returns true if the given page index points to the current page.
func (p *Paginator) IsActive(page int) bool {
	return p.Page() == page
}

// Offset method
// Returns the current offset.
func (p *Paginator) Offset() int {
	return (p.Page() - 1) * p.PerPageNums
}

// HasPages method
// Returns true if there is more than one page.
func (p *Paginator) HasPages() bool {
	return p.PageNums() > 1
}

// toInt64 helper
func toInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}
