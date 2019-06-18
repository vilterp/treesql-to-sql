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
