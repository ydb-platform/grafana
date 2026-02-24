CREATE TABLE IF NOT EXISTS seed_assignment_temp (
  id SERIAL,
  builtin_role TEXT,
  action TEXT,
  scope TEXT,
  role_name TEXT,
  PRIMARY KEY (id)
);
