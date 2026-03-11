UPSERT INTO folder (
    uid,
    org_id,
    title,
    created,
    updated
)
SELECT
    COALESCE(uid, ""),
    org_id,
    title,
    created,
    updated
FROM
    dashboard
WHERE
    is_folder == 1
;
