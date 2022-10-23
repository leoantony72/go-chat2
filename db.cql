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

CREATE TABLE chat.private_chat(
    id        VARCHAR PRIMARY KEY,
    sender    VARCHAR,
    receiver  VARCHAR ,
    msg       TEXT,
    timestamp timestamp
);
CREATE INDEX tb_private_chat_TIMESTAMP ON chat.private_chat(timestamp); 
CREATE INDEX tb_private_chat_SENDER ON chat.private_chat(sender); 

CREATE TABLE chat.group_chat(
    id        VARCHAR PRIMARY KEY,
    sender    VARCHAR,
    msg       TEXT,
    group     VARCHAR,
    timestamp timestamp
);
CREATE INDEX tb_group_chat_TIMSTAMP ON chat.group_chat(timestamp); 
CREATE INDEX tb_group_chat_SENDER ON chat.group_chat(sender); 