-- Used for testing purposes
-- drops tables if they already exist making it easier to rerun script later
ALTER TABLE messages 
DROP FOREIGN KEY member_id;
DROP TABLE IF EXISTS messages;

ALTER TABLE members 
DROP FOREIGN KEY room_id,
DROP FOREIGN KEY user_id;
DROP TABLE IF EXISTS members;

DROP TABLE IF EXISTS chatrooms;

DROP TABLE IF EXISTS users;


-- creates 
CREATE TABLE users(
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE chatrooms(
    room_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL, 
    password VARCHAR(255) NOT NULL
);

create TABLE members(
    member_id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT,
    user_id INT,
    FOREIGN KEY(room_id) REFERENCES chatrooms(room_id),
    FOREIGN KEY(user_id) REFERENCES users(user_id)
);

create TABLE messages(
    message_id INT AUTO_INCREMENT PRIMARY KEY,
    member_id INT,
    message VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(member_id) REFERENCES members(member_id)
);


-- inserts 
-- can be modified to insert more complex values
INSERT INTO users
    (username, email, password)
VALUES
    ('Vertigo', 'email@gmail.com', 'password'),
    ('Diva', 'diva@gmail.com', 'password'),
    ('Santa', 'santa@gmail.com', 'password');

SET @user_id_1 = LAST_INSERT_ID();
SET @user_id_2 = LAST_INSERT_ID();

INSERT INTO chatrooms
    (name, password)
VALUES
    ('The Room', 'password'),
    ('The Chat', 'password');

SET @chatroom_id = LAST_INSERT_ID();

INSERT INTO members
    (room_id,user_id)
VALUES
    (@chatroom_id, @user_id_2),
    (@chatroom_id, @user_id_1);

set @member_id = LAST_INSERT_ID();

INSERT INTO messages
    (member_id,message)
VALUES
    (@member_id, 'this');


