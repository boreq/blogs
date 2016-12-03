package database

var createUserSQL = `
CREATE TABLE "user" (
	id INTEGER PRIMARY KEY,

	username VARCHAR(1000) NOT NULL,
	password VARCHAR(1000) NOT NULL,

	UNIQUE(username)
)
`
var createUserSessionSQL = `
CREATE TABLE "user_session" (
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,

	key VARCHAR(1000) NOT NULL,
	last_seen DATETIME NOT NULL,

	FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(key)
)
`

var createBlogSQL = `
CREATE TABLE "blog" (
	id INTEGER PRIMARY KEY,

	internal_id INTEGER NOT NULL,
	title VARCHAR(1000) NOT NULL,

	UNIQUE(internal_id)
)
`

var createCategorySQL = `
CREATE TABLE "category" (
	id INTEGER PRIMARY KEY,
	blog_id INTEGER NOT NULL,

	name VARCHAR(1000) NOT NULL,

	FOREIGN KEY(blog_id) REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(blog_id, name)
)
`

var createPostSQL = `
CREATE TABLE "post" (
	id INTEGER PRIMARY KEY,
	category_id INTEGER NOT NULL,

	internal_id VARCHAR(1000) NOT NULL,
	title VARCHAR(1000) NOT NULL,
	summary VARCHAR(3000) NOT NULL,
	date DATETIME NOT NULL,

	FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(category_id, internal_id)
)
`

var createTagSQL = `
CREATE TABLE "tag" (
	id INTEGER PRIMARY KEY,

	name VARCHAR(1000) NOT NULL,

	UNIQUE(name)
)
`

var createPostToTagSQL = `
CREATE TABLE "post_to_tag" (
	id INTEGER PRIMARY KEY,
	post_id INTEGER NOT NULL,
	tag_id INTEGER NOT NULL,

	FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(tag_id) REFERENCES tag(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(post_id, tag_id)
)
`

var createUpdateSQL = `
CREATE TABLE "update" (
	id INTEGER PRIMARY KEY,
	blog_id INTEGER NOT NULL,

	started DATETIME NOT NULL,
	ended DATETIME NOT NULL,
	succeeded BOOLEAN NOT NULL,
	data TEXT,

	FOREIGN KEY(blog_id) REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE
)
`

var createSubscriptionSQL = `
CREATE TABLE "subscription" (
	id INTEGER PRIMARY KEY,
	blog_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,

	FOREIGN KEY(blog_id) REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(blog_id, user_id)
)
`
