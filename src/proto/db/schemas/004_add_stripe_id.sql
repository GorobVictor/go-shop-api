ALTER TABLE receipts
ADD COLUMN stripe_id VARCHAR(255),
ADD COLUMN stripe_status stripe_status NOT NULL DEFAULT 'pending';
