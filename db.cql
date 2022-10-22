CREATE KEYSPACE chat
WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};


CREATE TABLE chat.users(
    id VARCHAR ,
    username VARCHAR PRIMARY KEY
);

-- CREATE INDEX tb_users_id ON chat.users (id);

CREATE TABLE chat.room(
    id VARCHAR ,
    room_name VARCHAR PRIMARY KEY
);

-- CREATE INDEX tb_room_id ON chat.room (id);

CREATE TABLE chat.room_members(
    room_name VARCHAR,
    username VARCHAR,
    PRIMARY KEY(room_name,username)
);
-- CREATE INDEX tb_room_members_roomid ON chat.room_members (room_name);

CREATE TABLE chat.user_mapping(
    username    VARCHAR PRIMARY KEY,
    server_id   VARCHAR
);
-- CREATE INDEX tb_user_mapping_userid ON chat.user_mapping (username);

CREATE TABLE chat.message(
    sender    VARCHAR PRIMARY KEY,
    receiver  VARCHAR ,
    msg       TEXT,
    isgroup   BOOLEAN,
    group     VARCHAR  
);
-- CREATE INDEX tb_message_receiver ON chat.room_members (room_name);


