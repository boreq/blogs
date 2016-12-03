package core

import (
	"database/sql"
	"errors"
	"github.com/boreq/blogs/database"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/templates"
	"github.com/boreq/blogs/utils"
	verrors "github.com/boreq/blogs/views/errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type scannableTime struct {
	time.Time
}

func (t *scannableTime) Scan(src interface{}) error {
	updated, ok := src.([]uint8)
	if !ok {
		return errors.New("Invalid type, this is not []uint8")
	}
	tmp, err := time.Parse("2006-01-02 15:04:05-07:00", string(updated))
	if err != nil {
		return err
	}
	t.Time = tmp
	return nil
}

func (t scannableTime) String() string {
	return utils.ISO8601(t.Time)
}

type BlogResult struct {
	database.Blog
	SubscriptionID sql.NullInt64
	Updated        scannableTime
}

func blogs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get the data
	user_id := -1
	ctx := context.Get(r)
	if ctx.User.IsAuthenticated() {
		user_id = int(ctx.User.GetUser().ID)
	}

	var blogs = make([]BlogResult, 0)
	err := database.DB.Select(&blogs, `
		SELECT blog.*, subscription.id as subscription_id, MAX(post.date) AS updated
		FROM blog
		JOIN category ON category.blog_id = blog.id
		JOIN post ON post.category_id = category.id
		LEFT JOIN subscription ON subscription.blog_id = blog.id AND subscription.user_id=$1
		GROUP BY blog.id
		ORDER BY blog.title`, user_id)
	if err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}

	var data = templates.GetDefaultData(r)
	data["blogs"] = blogs
	if err := templates.RenderTemplateSafe(w, "core/blogs.tmpl", data); err != nil {
		verrors.InternalServerErrorWithStack(w, r, err)
		return
	}
}
