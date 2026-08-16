-- plaza_id_map moves out of JSON and into tables.
--
-- The shape stopped being a document that is only ever read and written whole.
-- It is a site with attributes and a list of facilities, and two rules the
-- JSON could not carry:
--
--   * a decode failure emptied every group silently, because the whole column
--     is one value - a column type mismatch surfaces as an error instead;
--   * the same facility could sit in two sites of one config, leaving the
--     lookup to pick whichever it scanned first, silently deciding which
--     providerId and clientId went to S&B.
--
-- Uniqueness is per config, not global: two different third parties serving
-- the same facility is exactly what the vendor fan-out exists for.
--
-- plaza_id_map itself is left in place and still populated. Nothing reads it
-- after this, and dropping it is a separate step - keeping it means the down
-- migration is just a DROP.
CREATE TABLE IF NOT EXISTS integrator_site
(
    id                   uuid primary key   default uuid_generate_v4(),
    integrator_config_id uuid      not null references integrator_config (id) on delete cascade,
    provider_id          int,
    client_id            varchar,
    vendor_location_id   varchar,
    created_at           timestamp not null default NOW(),
    updated_at           timestamp not null default NOW(),
    -- Redundant against the primary key, but it gives the facility table a
    -- composite foreign key to point at, which is what stops a facility row
    -- from naming a config its site does not belong to.
    UNIQUE (id, integrator_config_id)
);

CREATE TABLE IF NOT EXISTS integrator_site_facility
(
    site_id              uuid    not null,
    integrator_config_id uuid    not null,
    facility             varchar not null,
    PRIMARY KEY (site_id, facility),
    -- The rule the JSON could not enforce.
    UNIQUE (integrator_config_id, facility),
    FOREIGN KEY (site_id, integrator_config_id)
        REFERENCES integrator_site (id, integrator_config_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS integrator_site_config_idx
    ON integrator_site (integrator_config_id);

CREATE TRIGGER set_integrator_site_timestamp
    BEFORE UPDATE
    ON integrator_site
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Backfill from the grouped JSON.
--
-- MATERIALIZED so uuid_generate_v4() is evaluated once per group; without it
-- the two inserts below could see different ids for the same site.
WITH g AS MATERIALIZED (SELECT uuid_generate_v4()                  AS site_id,
                               c.id                                AS config_id,
                               grp.ordinality                      AS grp_no,
                               (grp.value ->> 'providerId')::int   AS provider_id,
                               grp.value ->> 'clientId'            AS client_id,
                               grp.value ->> 'vendorLocationId'    AS vendor_location_id,
                               grp.value -> 'facilities'           AS facilities
                        FROM integrator_config c,
                             jsonb_array_elements(c.plaza_id_map::jsonb -> 'groups')
                                 WITH ORDINALITY grp
                        WHERE c.plaza_id_map IS NOT NULL
                          AND jsonb_typeof(c.plaza_id_map::jsonb) = 'object'
                          AND jsonb_exists(c.plaza_id_map::jsonb, 'groups')),
     ins_site AS (
         INSERT INTO integrator_site (id, integrator_config_id, provider_id, client_id, vendor_location_id)
             SELECT site_id, config_id, provider_id, client_id, vendor_location_id
             FROM g)
INSERT
INTO integrator_site_facility (site_id, integrator_config_id, facility)
    -- DISTINCT ON keeps the first group that claims a facility, matching the
    -- old lookup, which returned the first match. Without it, legacy data
    -- carrying the same facility in two groups would fail the new unique
    -- constraint and abort the migration.
SELECT DISTINCT ON (g.config_id, f.value #>> '{}') g.site_id,
                                                   g.config_id,
                                                   f.value #>> '{}'
FROM g,
     jsonb_array_elements(g.facilities) f
ORDER BY g.config_id, f.value #>> '{}', g.grp_no;

-- A site whose facilities all lost the DISTINCT ON tie-break above would be
-- left with none, which the UI shows as an unnamed empty row.
DELETE
FROM integrator_site s
WHERE NOT EXISTS (SELECT 1
                  FROM integrator_site_facility f
                  WHERE f.site_id = s.id);
