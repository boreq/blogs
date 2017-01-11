package core

import (
	"github.com/boreq/blogs/database"
)

func getRecommendedBlogs(userId uint, limit uint) ([]popularBlogsResult, error) {
	var recommendedBlogs []popularBlogsResult
	if err := database.DB.Select(&recommendedBlogs,
		`SELECT b.*, MAX(p.date) AS updated, sa.id AS subscribed, COUNT(sb.id) AS score
		FROM blog b
		JOIN category c ON c.blog_id=b.id
		JOIN post p ON p.category_id=c.id
		LEFT JOIN subscription sa ON sa.blog_id=b.id AND sa.user_id=$1
		JOIN subscription sb ON sb.blog_id=b.id
		WHERE sb.user_id IN (
			SELECT DISTINCT s.user_id
			FROM subscription s
			WHERE s.blog_id IN (
				SELECT s.blog_id
				FROM subscription s
				JOIN "user" u ON u.id=s.user_id
				WHERE u.id=$1
			)
		)
		AND sa.id IS NULL
		GROUP BY b.id, sa.id, c.id
		ORDER BY score DESC
		LIMIT 5`, userId); err != nil {
		return nil, err
	}
	return recommendedBlogs, nil
}

func getRecommendedPosts(userId uint, limit uint) ([]popularPostsResult, error) {
	var recommendedPosts []popularPostsResult
	if err := database.DB.Select(&recommendedPosts,
		`SELECT p.*, c.*, b.*, sa.id AS starred, COUNT(sb.id) AS score
		FROM post p
		JOIN category c ON c.id=p.category_id
		JOIN blog b ON b.id=c.blog_id
		LEFT JOIN star sa ON sa.post_id=p.id AND sa.user_id=$1
		JOIN star sb ON sb.post_id=p.id
		WHERE sb.user_id IN (
			SELECT DISTINCT s.user_id
			FROM star s
			WHERE s.post_id IN (
				SELECT s.post_id
				FROM star s
				JOIN "user" u ON u.id=s.user_id
				WHERE u.id=$1
			)
		)
		AND sa.id IS NULL
		GROUP BY p.id, sa.id, c.id, b.id
		ORDER BY score DESC
		LIMIT 5`, userId); err != nil {
		return nil, err
	}
	return recommendedPosts, nil
}
