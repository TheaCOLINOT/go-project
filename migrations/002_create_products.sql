CREATE TABLE categories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(150) NOT NULL UNIQUE,
    description text,
    CONSTRAINT categories_name_not_empty CHECK (length(trim(name)) > 0)
);

CREATE TABLE products (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id varchar(64) NOT NULL UNIQUE,
    category_id uuid NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name varchar(200) NOT NULL,
    description text,
    price_ht_cents bigint NOT NULL,
    vat_rate_basis_points integer NOT NULL,
    stock integer NOT NULL DEFAULT 0,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT products_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT products_price_non_negative CHECK (price_ht_cents >= 0),
    CONSTRAINT products_vat_non_negative CHECK (vat_rate_basis_points >= 0),
    CONSTRAINT products_stock_non_negative CHECK (stock >= 0)
);

CREATE INDEX products_category_id_idx ON products (category_id);
CREATE INDEX products_is_active_idx ON products (is_active);
