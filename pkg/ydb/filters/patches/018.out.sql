UPSERT INTO folder (uid, org_id, title, created, updated)
SELECT
  g.folder_uid,
  g.org_id,
  g.title,
  g.created,
  g.updated
FROM
  (
    SELECT
      folder_uid,
      org_id,
      MAX(title) AS title,
      MAX(created) AS created,
      MAX(updated) AS updated
    FROM
      dashboard
    WHERE
      is_folder
    GROUP BY
      COALESCE(uid, "") AS folder_uid,
      org_id
  ) AS g;
