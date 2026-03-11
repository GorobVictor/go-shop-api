CREATE TYPE stripe_status AS ENUM(
    'pending',
    'succeeded',
    'failed',
    'canceled',
    'refunded'
);