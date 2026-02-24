UPDATE alert_rule_version
		SET rule_guid = alert_rule.guid
		FROM alert_rule
		WHERE alert_rule.uid = alert_rule_version.rule_uid
		  AND alert_rule.org_id = alert_rule_version.rule_org_id;
