CREATE TABLE feedbacks (
    id SERIAL PRIMARY KEY,
    xid UUID NOT NULL UNIQUE,
    user_id INTEGER NOT NULL,
    certificate_id INTEGER NOT NULL,
    result_clarity INTEGER,
    communication INTEGER,
    customer_service INTEGER,
    reccomend_service INTEGER,
    result_accuracy INTEGER,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (certificate_id) REFERENCES certificates(id) ON DELETE CASCADE
);
