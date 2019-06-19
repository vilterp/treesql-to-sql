CREATE TABLE posts (
  id INT PRIMARY KEY,
  title TEXT,
  body TEXT
);

CREATE TABLE comments (
  id INT PRIMARY KEY,
  post_id INT REFERENCES posts,
  body TEXT
);

INSERT INTO posts VALUES
  (1, 'hello world', 'bloop doop'),
  (2, 'hello again', 'woop loop');

INSERT INTO comments VALUES
  (1, 1, 'a comment'),
  (2, 1, 'another comment'),
  (3, 2, 'a third comment');
