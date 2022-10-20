CREATE KEYSPACE chat
WITH replication = {'class':'SimpleStrategy', 'replication_factor' : 1};


CREATE TABLE chat.users(
    id VARCHAR PRIMARY KEY,
    username VARCHAR
);

-- CREATE INDEX tb_users_id ON chat.users (id);

CREATE TABLE chat.room(
    id VARCHAR PRIMARY KEY,
    room_name VARCHAR
);

-- CREATE INDEX tb_room_id ON chat.room (id);

CREATE TABLE chat.room_members(
    room_id VARCHAR PRIMARY KEY,
    user_id VARCHAR
);
-- CREATE INDEX tb_room_members_roomid ON chat.room_members (room_id);

CREATE TABLE chat.user_mapping(
    user_id     VARCHAR PRIMARY KEY,
    server_id   VARCHAR
);
-- CREATE INDEX tb_user_mapping_userid ON chat.user_mapping (user_id);

CREATE TABLE chat.message(
    id        VARCHAR PRIMARY KEY,
    sender    VARCHAR,
    receiver  VARCHAR ,
    msg       TEXT,
    isgroup   BOOLEAN,
    group     VARCHAR  
);
-- CREATE INDEX tb_message_receiver ON chat.room_members (room_id);


