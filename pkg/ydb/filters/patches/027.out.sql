UPDATE `temp_user` SET created = ?, updated = ? WHERE created = 0 AND status in ('SignUpStarted', 'InvitePending')
