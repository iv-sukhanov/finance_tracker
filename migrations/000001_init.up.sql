CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table users (
    guid UUID not null default uuid_generate_v4() primary key,
    username VARCHAR(255) not null,
    telegram_id VARCHAR(255) not null,
    updated_at TIMESTAMP without time zone not null default now(),
    created_at TIMESTAMP without time zone not null default now()
);

create table spending_types (
    guid UUID not null default uuid_generate_v4() primary key,
    user_guid UUID references users (guid),
    s_type VARCHAR(255) not null,
    desctiption TEXT,
    amount int default 0,
    updated_at TIMESTAMP without time zone not null default now(),
    created_at TIMESTAMP without time zone not null default now()
);

create table spending_records (
    guid UUID not null default uuid_generate_v4() primary key,
    type_guid UUID references spending_types (guid),
    amount int default 0,
    desctiption TEXT,
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
    BEFORE UPDATE ON spending_types
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_notifications_modtime
    BEFORE UPDATE ON spending_records
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();
