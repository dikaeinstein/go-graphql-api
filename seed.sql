CREATE TABLE users (
  id serial PRIMARY KEY,
  name VARCHAR (50) NOT NULL,
  email VARCHAR (255) NOT NULL,
  age INT NOT NULL,
  profession VARCHAR (50) NOT NULL,
  friendly BOOLEAN NOT NULL
);

INSERT INTO users VALUES
  (1, 'kevin', 'kevin@email.com', 35, 'waiter', true),
  (2, 'angela', 'angela@email.com', 21, 'concierge', true),
  (3, 'alex', 'alex@email.com', 26, 'zoo keeper', false),
  (4, 'becky', 'becky@email.com', 67, 'retired', false),
  (5, 'kevin', 'kevin2@email.com', 15, 'in school', true),
  (6, 'frankie', 'frankie@email.com', 45, 'teller', true);
