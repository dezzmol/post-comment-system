CREATE TABLE IF NOT EXISTS users
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts
(
    id             SERIAL PRIMARY KEY,
    title          TEXT NOT NULL,
    content        TEXT NOT NULL,
    author_id      INTEGER,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    allow_comments BOOLEAN   DEFAULT TRUE,
    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS comments
(
    id         SERIAL PRIMARY KEY,
    post_id    INTEGER NOT NULL,
    text       TEXT    NOT NULL,
    author_id  INTEGER,
    reply_to   INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE SET NULL,
    FOREIGN KEY (reply_to) REFERENCES comments (id) ON DELETE SET NULL
);
