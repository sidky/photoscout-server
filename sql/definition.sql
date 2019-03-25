CREATE TABLE bookmarks (
  user_id VARCHAR(40) NOT NULL,
  photo_id VARCHAR(80) NOT NULL,
  updated_at TIMESTAMP DEFAULT now(),
  label VARCHAR(30),
  PRIMARY KEY (user_id, photo_id)
)