UPSERT INTO data_keys SELECT id AS name, keys.* WITHOUT name FROM data_keys AS keys WHERE id != name;
DELETE FROM data_keys WHERE name != id;
