BEGIN;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  telegram_user_id INT NOT NULL
);

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  telegram_user_id INT NOT NULL
  name CHAR(255)
  weight REAL(10)
  alreadyused BOOLEAN DEFAULT false
  inlist BOOLEAN DEFAULT false
  infridge BOOLEAN DEFAULT false
  intrash BOOLEAN DEFAULT false
  created_at DATE NOT NULL
  finished_at DATE DEFAULT NULL
  timerenable BOOLEAN DEFAULT false
);


COMMIT;