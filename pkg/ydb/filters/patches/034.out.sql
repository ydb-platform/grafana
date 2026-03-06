SELECT res.uid, res.is_folder, res.org_id
FROM (SELECT dashboard.id AS id, dashboard.uid AS uid, dashboard.is_folder AS is_folder, dashboard.org_id AS org_id, count(dashboard_acl.id) as count
      FROM dashboard
        LEFT JOIN dashboard_acl ON dashboard.id = dashboard_acl.dashboard_id
      WHERE dashboard.has_acl
      GROUP BY dashboard.id, dashboard.uid, dashboard.is_folder, dashboard.org_id, dashboard_acl.id) as res
WHERE res.count = 0
