CREATE TABLE "users"(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    image_path TEXT,
    email VARCHAR(254) UNIQUE NOT NULL,
    username VARCHAR(20) UNIQUE NOT NULL,
    hashed_password TEXT
);
CREATE TABLE user_admin(
    id UUID PRIMARY KEY REFERENCES "users"(id) ON DELETE CASCADE
);
CREATE TABLE category(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    activated BOOLEAN NOT NULL DEFAULT TRUE
); 
CREATE TABLE user_category(
    user_id UUID REFERENCES "users"(id) ON DELETE CASCADE,
    category_id UUID REFERENCES category(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id,category_id)
);

