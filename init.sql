CREATE TYPE sex AS ENUM ('man', 'woman');

CREATE TABLE orders(
  order_id SERIAL PRIMARY KEY,
  customer_id VARCHAR(255) NOT NULL,
  order_date DATE NOT NULL
);