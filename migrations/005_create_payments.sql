CREATE TABLE payments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    status varchar(32) NOT NULL DEFAULT 'pending',
    method varchar(64) NOT NULL,
    amount_cents bigint NOT NULL,
    simulated_reference varchar(255),
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT payments_status_not_empty CHECK (length(trim(status)) > 0),
    CONSTRAINT payments_method_not_empty CHECK (length(trim(method)) > 0),
    CONSTRAINT payments_amount_non_negative CHECK (amount_cents >= 0)
);

CREATE INDEX payments_order_id_idx ON payments (order_id);
CREATE INDEX payments_status_idx ON payments (status);
