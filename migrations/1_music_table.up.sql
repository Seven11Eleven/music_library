CREATE TABLE music(
    id SERIAL PRIMARY KEY,
    release_date DATE,
    title VARCHAR(255) NOT NULL,
    group_name VARCHAR(255),
    link TEXT
);

