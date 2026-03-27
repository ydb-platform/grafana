SELECT
  dashboard.id,
  dashboard.org_id,
  dashboard.uid,
  dashboard.title,
  dashboard.slug,
  dashboard_tag.term,
  dashboard.is_folder,
  dashboard.folder_id,
  dashboard.deleted,
  folder.uid AS folder_uid,
  folder.title AS folder_slug,
  folder.title AS folder_title,
  dashboard.title AS title
FROM (
  SELECT dashboard.id AS id, dashboard.title AS title
  FROM dashboard
  WHERE dashboard.org_id = $1
    AND dashboard.uid IN ($2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60)
    AND dashboard.is_folder = true
    AND dashboard.folder_uid IS NULL
    AND dashboard.uid != $61
    AND (dashboard.folder_uid != $62 OR dashboard.folder_uid IS NULL)
  ORDER BY title ASC
  LIMIT 50 OFFSET 0
) AS ids
INNER JOIN dashboard ON ids.id = dashboard.id
LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
ORDER BY dashboard.title ASC
