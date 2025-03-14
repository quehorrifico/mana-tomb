CREATE TABLE IF NOT EXISTS bulk_data (
    id UUID PRIMARY KEY,
    uri TEXT NOT NULL,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    download_uri TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    content_encoding TEXT NOT NULL
);