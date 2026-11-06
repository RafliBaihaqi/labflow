CREATE TABLE quotations (
    id SERIAL PRIMARY KEY,
    xid UUID NOT NULL UNIQUE,
    email TEXT NOT NULL,
    company_name TEXT NOT NULL,
    phone_number TEXT,
    message TEXT,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
