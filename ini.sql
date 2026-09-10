CREATE TABLE if NOT EXISTS Users(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    login text not null,
    password text not null
);

CREATE TABLE if NOT EXISTS Files(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);