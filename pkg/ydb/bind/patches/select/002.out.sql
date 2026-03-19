SELECT
    SUM(users) AS users,
    SUM(datasources) AS datasources,
    SUM(stars) AS stars,
    SUM(playlists) AS playlists,
    SUM(alerts) AS alerts,
    SUM(correlations) AS correlations,
    SUM(active_users) AS active_users,
    SUM(daily_active_users) AS daily_active_users,
    SUM(monthly_active_users) AS monthly_active_users,
    SUM(data_keys) AS data_keys,
    SUM(active_data_keys) AS active_data_keys,
    SUM(public_dashboards) AS public_dashboards,
    SUM(admins) AS admins,
    SUM(editors) AS editors,
    SUM(viewers) AS viewers,
    SUM(active_admins) AS active_admins,
    SUM(active_editors) AS active_editors,
    SUM(active_viewers) AS active_viewers,
    SUM(daily_active_admins) AS daily_active_admins,
    SUM(daily_active_editors) AS daily_active_editors,
    SUM(daily_active_viewers) AS daily_active_viewers
FROM (
    (
        SELECT
            COUNT(*) AS users
        FROM
            `user`
        WHERE
            is_service_account == false
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS datasources
        FROM
            `data_source`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS stars
        FROM
            `star`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS playlists
        FROM
            `playlist`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS alerts
        FROM
            `alert`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS correlations
        FROM
            `correlation`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS active_users
        FROM
            `user`
        WHERE
            is_service_account == false AND last_seen_at > ?
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS daily_active_users
        FROM
            `user`
        WHERE
            is_service_account == false AND last_seen_at > ?
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS monthly_active_users
        FROM
            `user`
        WHERE
            is_service_account == false AND last_seen_at > ?
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS provisioned_dashboards
        FROM
            `dashboard_provisioning`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS snapshots
        FROM
            `dashboard_snapshot`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS dashboard_versions
        FROM
            `dashboard_version`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS annotations
        FROM
            `annotation`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS teams
        FROM
            `team`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS auth_tokens
        FROM
            `user_auth_token`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS alert_rules
        FROM
            `alert_rule`
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS api_keys
        FROM
            `api_key`
        WHERE
            service_account_id IS NULL
    )
    UNION ALL
    (
        SELECT
            COUNT(id) AS library_panels
        FROM
            `library_element`
        WHERE
            kind == ?
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS data_keys
        FROM
            `data_keys`
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS active_data_keys
        FROM
            `data_keys`
        WHERE
            active == true
    )
    UNION ALL
    (
        SELECT
            COUNT(*) AS public_dashboards
        FROM
            `dashboard_public`
    )
    UNION ALL
    (
        SELECT
            MIN(timestamp) AS database_created_time
        FROM
            `migration_log`
    )
    UNION ALL
    (
        SELECT
            COUNT(DISTINCT (`rule_group`)) AS rule_groups
        FROM
            `alert_rule`
    )
    UNION ALL
    (
        SELECT
            0 AS admins
    )
    UNION ALL
    (
        SELECT
            0 AS editors
    )
    UNION ALL
    (
        SELECT
            0 AS viewers
    )
    UNION ALL
    (
        SELECT
            0 AS active_admins
    )
    UNION ALL
    (
        SELECT
            0 AS active_editors
    )
    UNION ALL
    (
        SELECT
            0 AS active_viewers
    )
    UNION ALL
    (
        SELECT
            0 AS daily_active_admins
    )
    UNION ALL
    (
        SELECT
            0 AS daily_active_editors
    )
    UNION ALL
    (
        SELECT
            0 AS daily_active_viewers
    )
);
