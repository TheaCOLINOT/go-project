CREATE TABLE carts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id varchar(64) NOT NULL UNIQUE,
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE cart_items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id varchar(64) NOT NULL UNIQUE,
    cart_id uuid NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity integer NOT NULL,
    CONSTRAINT cart_items_quantity_positive CHECK (quantity > 0),
    CONSTRAINT cart_items_cart_product_unique UNIQUE (cart_id, product_id)
);

CREATE INDEX cart_items_product_id_idx ON cart_items (product_id);
