INSERT INTO folder (uid, org_id, title, created, updated)
SELECT * FROM (SELECT uid, org_id, title, created, updated FROM dashboard WHERE is_folder = 1) AS derived
ON DUPLICATE KEY UPDATE title=derived.title, updated=derived.updated
