CREATE TABLE orders (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    business_id varchar(64) NOT NULL UNIQUE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_by_admin_id uuid REFERENCES users(id) ON DELETE RESTRICT,
    status varchar(32) NOT NULL DEFAULT 'pending',
    total_ttc_cents bigint NOT NULL,
    cancellation_reason text,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT orders_status_not_empty CHECK (length(trim(status)) > 0),
    CONSTRAINT orders_total_non_negative CHECK (total_ttc_cents >= 0)
);

CREATE INDEX orders_user_id_idx ON orders (user_id);
CREATE INDEX orders_created_by_admin_id_idx ON orders (created_by_admin_id);
CREATE INDEX orders_status_idx ON orders (status);
CREATE INDEX orders_created_at_idx ON orders (created_at);

CREATE TABLE order_items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_name_snapshot varchar(200) NOT NULL,
    quantity integer NOT NULL,
    unit_price_ht_cents bigint NOT NULL,
    vat_rate_basis_points integer NOT NULL,
    unit_price_ttc_cents bigint NOT NULL,
    CONSTRAINT order_items_product_name_not_empty CHECK (length(trim(product_name_snapshot)) > 0),
    CONSTRAINT order_items_quantity_positive CHECK (quantity > 0),
    CONSTRAINT order_items_unit_price_ht_non_negative CHECK (unit_price_ht_cents >= 0),
    CONSTRAINT order_items_vat_non_negative CHECK (vat_rate_basis_points >= 0),
    CONSTRAINT order_items_unit_price_ttc_non_negative CHECK (unit_price_ttc_cents >= 0)
);

CREATE INDEX order_items_order_id_idx ON order_items (order_id);
CREATE INDEX order_items_product_id_idx ON order_items (product_id);

CREATE TABLE order_status_history (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    changed_by_user_id uuid REFERENCES users(id) ON DELETE SET NULL,
    old_status varchar(32),
    new_status varchar(32) NOT NULL,
    reason text,
    changed_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT order_status_history_new_status_not_empty CHECK (length(trim(new_status)) > 0)
);

CREATE INDEX order_status_history_order_id_idx ON order_status_history (order_id);
CREATE INDEX order_status_history_changed_by_user_id_idx ON order_status_history (changed_by_user_id);
