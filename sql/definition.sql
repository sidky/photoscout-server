CREATE TABLE starred_photo (
  user_id VARCHAR(40) PRIMARY KEY,
  photo_id VARCHAR(80) NOT NULL,
  updated_at TIMESTAMP DEFAULT now(),
  label VARCHAR(30)
)