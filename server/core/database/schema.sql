CREATE TABLE
    hamlet_users (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        display_name VARCHAR(64) NOT NULL,
        is_banned BOOLEAN DEFAULT FALSE NOT NULL,
        public_key BYTEA NOT NULL,
        is_verified BOOLEAN DEFAULT FALSE NOT NULL
    );

CREATE TABLE
    hamlet_user_permissions (
        user_id FOREIGN KEY REFERENCES hamlet_users.id NOT NULL,
        is_owner BOOLEAN DEFAULT FALSE NOT NULL,
        manage_admin BOOLEAN DEFAULT FALSE NOT NULL,
        manage_users BOOLEAN DEFAULT FALSE NOT NULL,
        manage_all_lodges BOOLEAN DEFAULT FALSE NOT NULL,
        manage_own_lodges BOOLEAN DEFAULT FALSE NOT NULL,
    );

CREATE TABLE
    hamlet_lodge (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        lodge_name VARCHAR(64) NOT NULL,
        lodge_description VARCHAR(256) NOT NULL DEFAULT '',
        discoverable BOOLEAN DEFAULT FALSE NOT NULL,
        vanity_url VARCHAR(32) DEFAULT NULL UNIQUE,
    );

CREATE TABLE
    hamlet_lodge_users (
        user_id uuid FOREIGN KEY REFERENCES hamlet_users.id NOT NULL,
        -- Keep record of users which have left the lodge
        is_member BOOLEAN DEFAULT TRUE NOT NULL,
        -- Keep record of banned users
        is_banned BOOLEAN DEFAULT FALSE NOT NULL,
        -- Only set if changed from hamlet wide user display name
        display_name VARCHAR(64) DEFAULT NULL,
    );

CREATE TABLE
    hamlet_lodge_rooms (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        provider_id VARCHAR(128) NOT NULL,
        room_type_id VARCHAR(64) NOT NULL,
        name VARCHAR(32) NOT NULL,
    );

CREATE TABLE
    hamlet_lodge_groups (
        id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
        name VARCHAR(64) NOT NULL,
        visible BOOLEAN DEFAULT TRUE,
    );

CREATE TABLE
    hamlet_lodge_groups_users (
        user_id uuid FOREIGN KEY REFERENCES hamlet_lodge_users.id NOT NULL,
        group_id uuid FOREIGN KEY REFERENCES hamlet_lodge_groups.id NOT NULL,
        PRIMARY KEY (user_id, group_id)
    );

CREATE TABLE
    hamlet_lodge_group_permissions (
        group_id uuid FOREIGN KEY REFERENCES hamlet_lodge_groups.id NOT NULL,
        permission_id VARCHAR(64) NOT NULL,
        -- Up to 128 chars so that URLs can be used for better unique-ness
        -- Provider IDs should not begin with "hamlet/" unless they are built in (e.g. voice, text, core)
        provider_id VARCHAR(128) NOT NULL,
        -- Only TRUE/ FALSE permissions are supported
        permission_value BOOLEAN DEFAULT FALSE NOT NULL,
        PRIMARY KEY (permission_id, provider_id, group_id)
    );

CREATE TABLE
    hamlet_auth_challenges (
        public_key BYTEA NOT NULL PRIMARY KEY,
        challenge_data BYTEA NOT NULL,
        expires DATETIME NOT NULL,
    );