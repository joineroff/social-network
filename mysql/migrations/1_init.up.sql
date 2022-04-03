CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    login varchar(320) NOT NULL UNIQUE,
    password varchar(60) NOT NULL,
    first_name varchar(50),
    last_name varchar(50),
    gender int,
    interests varchar(500),
    city varchar(60)
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id int NOT NULL,
    friend_id int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (friend_id) REFERENCES users (id)
);
