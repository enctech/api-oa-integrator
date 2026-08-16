-- plaza_id_map becomes a list of groups, one per site.
--
-- `name` is the config identity and is uniquely indexed, so several sites
-- sharing an integrator cannot be several rows - four sites all named TNG
-- collide on integrator_config_name_key. They become groups inside one row
-- instead, each owning the providerId S&B expects and the clientId the vendor
-- issued.
--
-- Two older layouts are converted:
--   {"1225":"ACL"}
--   {"1225":{"vendorLocationId":"ACL","clientId":"X","providerId":2}}
-- into:
--   {"groups":[{"providerId":..,"clientId":..,"vendorLocationId":"ACL",
--               "facilities":["1225"]}]}
--
-- Facilities are grouped by their effective (providerId, clientId,
-- vendorLocationId), falling back to the config-level columns where the entry
-- did not override them. Grouping on vendorLocationId too means a site whose
-- facilities pointed at different vendor locations splits into one group per
-- location rather than silently losing one of them.
-- Those columns are left in place for now; nothing reads them after this, and
-- dropping them is a separate step once the grouped shape has run in staging.
UPDATE integrator_config c
SET plaza_id_map = sub.grouped
FROM (SELECT c2.id,
             json_build_object('groups', json_agg(g.grp)) AS grouped
      FROM integrator_config c2,
           LATERAL (
               SELECT json_build_object(
                              'providerId', k.provider_id,
                              'clientId', k.client_id,
                              'vendorLocationId', k.vendor_location_id,
                              'facilities', json_agg(k.facility)
                      ) AS grp
               FROM (SELECT e.key                                    AS facility,
                            CASE jsonb_typeof(e.value)
                                WHEN 'string' THEN e.value #>> '{}'
                                ELSE e.value ->> 'vendorLocationId'
                                END                                  AS vendor_location_id,
                            COALESCE(
                                    NULLIF(CASE jsonb_typeof(e.value)
                                               WHEN 'object' THEN e.value ->> 'clientId'
                                               END, ''),
                                    c2.client_id)                    AS client_id,
                            COALESCE(
                                    NULLIF(CASE jsonb_typeof(e.value)
                                               WHEN 'object' THEN e.value ->> 'providerId'
                                               END, '')::int,
                                    c2.provider_id)                  AS provider_id
                     FROM jsonb_each(c2.plaza_id_map::jsonb) e) k
               GROUP BY k.provider_id, k.client_id, k.vendor_location_id
               ) g
      WHERE c2.plaza_id_map IS NOT NULL
        AND jsonb_typeof(c2.plaza_id_map::jsonb) = 'object'
        AND NOT jsonb_exists(c2.plaza_id_map::jsonb, 'groups')
      GROUP BY c2.id) sub
WHERE c.id = sub.id;

-- provider_id and client_id now live on the group, so the columns stop being
-- written. provider_id was NOT NULL, which would reject every new config.
ALTER TABLE integrator_config
    ALTER COLUMN provider_id DROP NOT NULL;

-- The unique index is not merely unused now, it is wrong: one server runs
-- several sites and each site's providerId is its own, so two configs may
-- legitimately carry the same number. This is the constraint the original
-- report ran into - "ProviderID boleh set 1 je utk 1 server".
DROP INDEX IF EXISTS integrator_config_provider_id_key;
