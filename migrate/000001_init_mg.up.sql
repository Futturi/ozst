CREATE TABLE IF NOT EXISTS posts (
                                     id bigserial PRIMARY KEY,
                                     title TEXT,
                                     content TEXT,
                                     commentsAllowed BOOLEAN
);

CREATE TABLE IF NOT EXISTS comments (
                                        id bigserial PRIMARY KEY,
                                        postId INTEGER,
                                        parentId INTEGER,
                                        content TEXT,
                                        FOREIGN KEY(postId) REFERENCES posts(id)
    );