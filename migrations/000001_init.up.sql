CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users (
    guid UUID not null default uuid_generate_v4() primary key,
    username VARCHAR(255) not null unique,
    telegram_id VARCHAR(255) not null unique,
    updated_at TIMESTAMP without time zone not null default now(),
    created_at TIMESTAMP without time zone not null default now()
);

create table spending_categories (
    guid UUID not null default uuid_generate_v4() primary key,
    user_guid UUID references users (guid),
    category VARCHAR(255) not null,
    description TEXT,
    amount NUMERIC(20, 0) default 0,
    updated_at TIMESTAMP without time zone not null default now(),
    created_at TIMESTAMP without time zone not null default now(),
    CONSTRAINT unique_columns UNIQUE (user_guid, category)
);

create table spending_records (
    guid UUID not null default uuid_generate_v4() primary key,
    category_guid UUID references spending_categories (guid),
    amount NUMERIC(10, 0) default 0,
    description TEXT,
    updated_at TIMESTAMP without time zone not null default now(),
    created_at TIMESTAMP without time zone not null default now()
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_filters_modtime
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();
    
CREATE TRIGGER update_notifications_modtime
    BEFORE UPDATE ON spending_categories
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_notifications_modtime
    BEFORE UPDATE ON spending_records
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();
