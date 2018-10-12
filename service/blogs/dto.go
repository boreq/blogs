package blogs

import (
	"github.com/boreq/blogs/dto"
)

type ListOut struct {
	Page  dto.PageOut   `json:"page"`
	Blogs []dto.BlogOut `json:"blogs"`
}
