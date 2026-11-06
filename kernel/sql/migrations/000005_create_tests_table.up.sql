CREATE TABLE tests (
    id SERIAL PRIMARY KEY,
    xid UUID NOT NULL UNIQUE,
    reservation_id INTEGER NOT NULL,
    test_list_id INTEGER NOT NULL,
    method VARCHAR(255) NOT NULL,
    criteria VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    test_unit VARCHAR(255) NOT NULL,
    test_result VARCHAR(255),
    test_summary VARCHAR(255),
    approved_date TIMESTAMP WITH TIME ZONE,
    notes VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (reservation_id) REFERENCES reservations(id) ON DELETE CASCADE,
    FOREIGN KEY (test_list_id) REFERENCES test_lists(id) ON DELETE CASCADE
);
