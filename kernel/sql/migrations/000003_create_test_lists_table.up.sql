CREATE TABLE test_lists (
    id SERIAL PRIMARY KEY,
    xid UUID NOT NULL,
    category_id INTEGER NOT NULL,
    value VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    sample_amount_liquid INTEGER,
    sample_amount_solid INTEGER,
    sample_amount_powder INTEGER,
    leadtime INTEGER,
    version INTEGER NOT NULL DEFAULT 1,
    is_current BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    UNIQUE(xid, version)
);
