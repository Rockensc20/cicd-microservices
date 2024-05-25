CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tag
(
    id SERIAL,
    name TEXT NOT NULL,
    CONSTRAINT tag_pkey PRIMARY KEY (id)
); 

CREATE TABLE IF NOT EXISTS productToTagAssignment
(
    id SERIAL,
    productID integer NOT NULL,
	tagID integer NOT NULL,
    CONSTRAINT productToTagAssignment_pkey PRIMARY KEY (id),
	CONSTRAINT product_fkey FOREIGN KEY (productID) REFERENCES products (id) ON DELETE CASCADE,
	CONSTRAINT tag_fkey FOREIGN KEY (tagID) REFERENCES tag (id) ON DELETE CASCADE
); 