DROP TABLE comments;
DROP TABLE posts;

CREATE TABLE users (
  id INT PRIMARY KEY,
  name TEXT
);

CREATE TABLE posts (
  id INT PRIMARY KEY,
  author_id INT REFERENCES users,
  title TEXT,
  body TEXT
);

CREATE TABLE comments (
  id INT PRIMARY KEY,
  author_id INT REFERENCES users,
  post_id INT REFERENCES posts,
  body TEXT
);

INSERT INTO users VALUES
  (1, 'Alice'),
  (2, 'Bob');

INSERT INTO posts VALUES
  (1, 1, 'hello world', 'bloop doop'),
  (2, 2, 'hello again', 'woop loop');

INSERT INTO comments VALUES
  (1, 1, 1, 'a comment'),
  (2, 2, 1, 'another comment'),
  (3, 2, 2, 'a third comment');
