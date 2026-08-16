-- plaza_id_map was left populated by the up migration, so the tables are the
-- only new state and dropping them is the whole rollback.
DROP TRIGGER IF EXISTS set_integrator_site_timestamp ON integrator_site;
DROP TABLE IF EXISTS integrator_site_facility;
DROP TABLE IF EXISTS integrator_site;
