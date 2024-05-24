-- CREATE TABLE POST
CREATE TABLE Post (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(10) CHECK (status IN ('draft', 'publish')) NOT NULL,
    publish_date TIMESTAMP,
    deleted_at TIMESTAMP
);

-- CREATE TABLE TAG
CREATE TABLE Tag (
    id SERIAL PRIMARY KEY,
    label VARCHAR(255) UNIQUE NOT NULL,
    deleted_at TIMESTAMP
);


-- CREATE TABLE POSTTAG
CREATE TABLE PostTag (
    post_id INT,
    tag_id INT,
    PRIMARY KEY (post_id, tag_id),
    FOREIGN KEY (post_id) REFERENCES Post(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES Tag(id) ON DELETE CASCADE,
    deleted_at TIMESTAMP
);

-- CREATE TABLE USERS
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    deleted_at TIMESTAMP
);