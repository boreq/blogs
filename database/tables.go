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
	subscriptions INT NOT NULL DEFAULT 0,

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
	stars INT NOT NULL DEFAULT 0,

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
	date    DATETIME NOT NULL,

	FOREIGN KEY(blog_id) REFERENCES blog(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(blog_id, user_id)
)
`

var createInsertSubscriptionTriggerSQL = `
CREATE TRIGGER update_subscriptions_insert AFTER INSERT ON subscription 
BEGIN
	UPDATE blog SET subscriptions=((SELECT subscriptions FROM blog WHERE blog.id=new.blog_id)+1) WHERE blog.id=new.blog_id;
END;
`

var createDeleteSubscriptionTriggerSQL = `
CREATE TRIGGER update_subscriptions_delete AFTER DELETE ON subscription 
BEGIN
	UPDATE blog SET subscriptions=((SELECT subscriptions FROM blog WHERE blog.id=old.blog_id)-1) WHERE blog.id=old.blog_id;
END;
`

var createStarSQL = `
CREATE TABLE "star" (
	id INTEGER PRIMARY KEY,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	date    DATETIME NOT NULL,

	FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE ON UPDATE CASCADE,
	FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE ON UPDATE CASCADE,
	UNIQUE(post_id, user_id)
)
`

var createInsertStarTriggerSQL = `
CREATE TRIGGER update_stars_insert AFTER INSERT ON star 
BEGIN
	UPDATE post SET stars=((SELECT stars FROM post WHERE post.id=new.post_id)+1) WHERE post.id=new.post_id;
END;
`

var createDeleteStarTriggerSQL = `
CREATE TRIGGER update_stars_delete AFTER DELETE ON star 
BEGIN
	UPDATE post SET stars=((SELECT stars FROM post WHERE post.id=old.post_id)-1) WHERE post.id=old.post_id;
END;
`

var triggersPostgreSQL = `
DROP FUNCTION IF EXISTS update_subscriptions_insert();
CREATE FUNCTION update_subscriptions_insert() RETURNS trigger AS $update_subscriptions_insert$
    BEGIN
	UPDATE blog SET subscriptions=((SELECT subscriptions FROM blog WHERE blog.id=NEW.blog_id)+1) WHERE blog.id=NEW.blog_id;
        RETURN NEW;
    END;
$update_subscriptions_insert$ LANGUAGE plpgsql;

CREATE TRIGGER update_subscriptions_insert AFTER INSERT ON subscription 
	FOR EACH ROW EXECUTE PROCEDURE update_subscriptions_insert();


DROP FUNCTION IF EXISTS update_subscriptions_delete();
CREATE FUNCTION update_subscriptions_delete() RETURNS trigger AS $update_subscriptions_delete$
    BEGIN
	UPDATE blog SET subscriptions=((SELECT subscriptions FROM blog WHERE blog.id=OLD.blog_id)-1) WHERE blog.id=OLD.blog_id;
        RETURN OLD;
    END;
$update_subscriptions_delete$ LANGUAGE plpgsql;

CREATE TRIGGER update_subscriptions_delete AFTER DELETE ON subscription 
	FOR EACH ROW EXECUTE PROCEDURE update_subscriptions_delete();


DROP FUNCTION IF EXISTS update_stars_insert();
CREATE FUNCTION update_stars_insert() RETURNS trigger AS $update_stars_insert$
    BEGIN
	UPDATE post SET stars=((SELECT stars FROM post WHERE post.id=NEW.post_id)+1) WHERE post.id=NEW.post_id;
        RETURN NEW;
    END;
$update_stars_insert$ LANGUAGE plpgsql;

CREATE TRIGGER update_stars_insert AFTER INSERT ON star 
	FOR EACH ROW EXECUTE PROCEDURE update_stars_insert();


DROP FUNCTION IF EXISTS update_stars_delete();
CREATE FUNCTION update_stars_delete() RETURNS trigger AS $update_stars_delete$
    BEGIN
	UPDATE post SET stars=((SELECT stars FROM post WHERE post.id=OLD.post_id)-1) WHERE post.id=OLD.post_id;
        RETURN OLD;
    END;
$update_stars_delete$ LANGUAGE plpgsql;

CREATE TRIGGER update_stars_delete AFTER DELETE ON star 
	FOR EACH ROW EXECUTE PROCEDURE update_stars_delete();
`
