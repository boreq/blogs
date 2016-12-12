package blogs

import (
	"github.com/boreq/blogs/blogs/loaders"
	"github.com/boreq/blogs/blogs/loaders/datarebellion"
	"github.com/boreq/blogs/blogs/loaders/eevee"
	"github.com/boreq/blogs/blogs/loaders/golang"
	"github.com/boreq/blogs/blogs/loaders/ilikebigbits"
	"github.com/boreq/blogs/blogs/loaders/lucumr"
	"github.com/boreq/blogs/blogs/loaders/yegor256"
)

// Blogs is a map mapping internal IDs of all supported blogs to their loaders.
var Blogs map[uint]loaders.Blog

func init() {
	Blogs = make(map[uint]loaders.Blog)
	Blogs[0] = eevee.New()
	Blogs[1] = lucumr.New()
	Blogs[2] = ilikebigbits.New()
	Blogs[4] = datarebellion.New()
	Blogs[5] = yegor256.New()
	Blogs[6] = golang.New()
}
