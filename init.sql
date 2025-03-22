CREATE TYPE sex AS ENUM ('man', 'woman');
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled');

CREATE TABLE orders (
  order_id SERIAL PRIMARY KEY,
  customer_id TEXT NOT NULL,
  status order_status NOT NULL DEFAULT 'pending',
  order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
