DROP INDEX login_email_idx;
DROP INDEX login_created_at_idx;
DROP INDEX login_updated_at_idx;

DROP TRIGGER login_set_updated_at_trg ON login;
DROP TABLE login;
