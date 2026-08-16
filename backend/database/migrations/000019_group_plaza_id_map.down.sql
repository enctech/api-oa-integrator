-- Flatten groups back to the per-facility map.
--
-- providerId and clientId survive as per-entry overrides, which the flat
-- format carries, so no group data is lost.
--
-- This can legitimately fail at the last step. The grouped format allows two
-- configs to share a providerId; the unique index recreated below forbids it.
-- If that happens the data genuinely does not fit the old shape, and the
-- failure is the correct outcome rather than something to force past.

-- Column values come back first: flattening overwrites plaza_id_map, so the
-- groups have to be read before they are gone.
UPDATE integrator_config c
SET provider_id = COALESCE(c.provider_id, (c.plaza_id_map::jsonb -> 'groups' -> 0 ->> 'providerId')::int),
    client_id   = COALESCE(c.client_id, c.plaza_id_map::jsonb -> 'groups' -> 0 ->> 'clientId')
WHERE c.plaza_id_map IS NOT NULL
  AND jsonb_exists(c.plaza_id_map::jsonb, 'groups');

UPDATE integrator_config c
SET plaza_id_map = sub.flat
FROM (SELECT c2.id,
             json_object_agg(f.facility,
                             json_build_object(
                                     'vendorLocationId', f.vendor_location_id,
                                     'clientId', f.client_id,
                                     'providerId', f.provider_id
                             )) AS flat
      FROM integrator_config c2,
           LATERAL (
               SELECT fac.value #>> '{}'                AS facility,
                      grp.value ->> 'vendorLocationId'  AS vendor_location_id,
                      grp.value ->> 'clientId'          AS client_id,
                      (grp.value ->> 'providerId')::int AS provider_id
               FROM jsonb_array_elements(c2.plaza_id_map::jsonb -> 'groups') grp,
                    jsonb_array_elements(grp.value -> 'facilities') fac
               ) f
      WHERE c2.plaza_id_map IS NOT NULL
        AND jsonb_exists(c2.plaza_id_map::jsonb, 'groups')
      GROUP BY c2.id) sub
WHERE c.id = sub.id;

UPDATE integrator_config SET provider_id = 0 WHERE provider_id IS NULL;

ALTER TABLE integrator_config
    ALTER COLUMN provider_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS integrator_config_provider_id_key
    ON integrator_config (provider_id)
    WHERE deleted_at IS NULL;
