package blogs

import (
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/blogs/loaders/eevee"
)

// Blogs is a map mapping internal IDs of all supported blogs to their loaders.
var Blogs map[uint]loaders.Blog

func init() {
	Blogs = make(map[uint]loaders.Blog)
	Blogs[0] = eevee.New()
}
