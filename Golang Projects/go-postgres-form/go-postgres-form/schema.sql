CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    class TEXT NOT NULL,
    colocation TEXT NOT NULL,
    description TEXT NOT NULL,
    hostname TEXT NOT NULL,
    mgmt_ip INET NOT NULL,
    serial_number TEXT NOT NULL,
    device_type TEXT NOT NULL
);
