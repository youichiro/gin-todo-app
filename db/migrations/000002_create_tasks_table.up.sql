CREATE TABLE IF NOT EXISTS tasks(
  id SERIAL PRIMARY KEY NOT NULL,
  title VARCHAR (50) NOT NULL,
  done BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER refresh_users_updated_at_step1
  BEFORE UPDATE ON tasks FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step1();
CREATE TRIGGER refresh_users_updated_at_step2
  BEFORE UPDATE OF updated_at ON tasks FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step2();
CREATE TRIGGER refresh_users_updated_at_step3
  BEFORE UPDATE ON tasks FOR EACH ROW
  EXECUTE PROCEDURE refresh_updated_at_step3();

