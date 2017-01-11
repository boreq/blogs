package updater

import (
	"github.com/boreq/blogs/database"
	"github.com/boreq/sqlx"
	"strings"
	"time"
)

var uniqueConstraintFailedErrorText = "UNIQUE constraint failed"

func getBlog(internalID uint) (database.Blog, error) {
	_, err := database.DB.NamedExec(
		"INSERT INTO blog (internal_id, title) VALUES (:internal_id, :title)",
		&database.Blog{InternalID: internalID, Title: ""})
	if err != nil && !isUniqueConstraintError(err) {
		log.Printf("getBlog, error: %s", err.Error())
	}
	blog := database.Blog{}
	err = database.DB.Get(&blog,
		"SELECT * FROM blog WHERE internal_id=$1 LIMIT 1",
		internalID)
	return blog, err
}

func getCategory(tx *sqlx.Tx, blog database.Blog, name string) (database.Category, error) {
	_, err := tx.Exec(
		`INSERT INTO category (blog_id, name)
		SELECT $1, $2
		WHERE NOT EXISTS (
		    SELECT 1
		    FROM category
		    WHERE blog_id=$3 AND name=$4
		)`,
		blog.ID, name, blog.ID, name)
	if err != nil && !isUniqueConstraintError(err) {
		log.Printf("getCategory, error: %s", err.Error())
	}
	category := database.Category{}
	err = tx.Get(&category,
		"SELECT * FROM category WHERE blog_id=$1 AND name=$2 LIMIT 1",
		blog.ID, name)
	return category, err
}

func getPost(tx *sqlx.Tx, category database.Category, internalID string) (database.Post, error) {
	_, err := tx.Exec(
		`INSERT INTO post (category_id, internal_id, title, summary, date)
		VALUES ($1, $2, $3, $4, $5)`,
		category.ID, internalID, "", "", time.Time{})
	if err != nil && !isUniqueConstraintError(err) {
		log.Printf("getPost, error: %s", err.Error())
	}
	post := database.Post{}
	err = tx.Get(&post,
		"SELECT * FROM post WHERE category_id=$1 AND internal_id=$2 LIMIT 1",
		category.ID, internalID)
	return post, err
}

func getTags(post database.Post) ([]database.Tag, error) {
	var tags []database.Tag
	err := database.DB.Select(&tags,
		`SELECT tag.*
		FROM tag
		JOIN post_to_tag ON post_to_tag.tag_id=tag.id
		JOIN post ON post.id=post_to_tag.post_id
		WHERE post.id=$1`, post.ID)
	return tags, err
}

func getTag(name string) (database.Tag, error) {
	_, err := database.DB.Exec(`
		INSERT INTO tag (name)
		SELECT $1
		WHERE NOT EXISTS (
		    SELECT 1
		    FROM tag
		    WHERE name=$2
		)
	`, name, name)

	if err != nil && !isUniqueConstraintError(err) {
		log.Printf("getTag, error: %s", err.Error())
	}
	tag := database.Tag{}
	err = database.DB.Get(&tag, "SELECT * FROM tag WHERE name=$1 LIMIT 1", name)
	return tag, err
}

type postWithCategory struct {
	database.Post
	database.Category
}

func getPostsForDeleteOldPosts(blog database.Blog) ([]postWithCategory, error) {
	var posts []postWithCategory
	err := database.DB.Select(&posts,
		`SELECT post.*, category.*
		FROM post
		JOIN category ON post.category_id=category.id
		JOIN blog ON category.blog_id=blog.id
		WHERE blog.id=$1`, blog.ID)
	return posts, err
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), uniqueConstraintFailedErrorText)
}
