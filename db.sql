DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS tweets;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    followers_id JSON,
    following_id JSON,
    feed JSON
);

CREATE TABLE tweets (
    id BIGINT AUTO_INCREMENT PRIMARY KEY, 
    message TEXT NOT NULL,
    timestamp DATETIME NOT NULL,
    author_id BIGINT NOT NULL,
    FOREIGN KEY (author_id) REFERENCES users(id)
);

CREATE TABLE follows (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    follower_id BIGINT NOT NULL,
    followed_id BIGINT NOT NULL,
    FOREIGN KEY (follower_id) REFERENCES users(id),
    FOREIGN KEY (followed_id) REFERENCES users(id),
    UNIQUE (follower_id, followed_id)
);
