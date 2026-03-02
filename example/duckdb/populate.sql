DROP TABLE IF EXISTS collections;
CREATE TABLE collections (id TEXT, content JSON);
INSERT INTO collections
SELECT id
    , {
        id: id,
        "type": 'Collection',
        stac_version: stac_version,
        description: description,
        license: license,
        extent: struct_update(
            extent,
            temporal := struct_update(
                extent.temporal,
                "interval" := [
                    list_transform(
                        extent.temporal.interval[1],
                        lambda x: strftime(x, '%Y-%m-%dT%H:%M:%S.%fZ')
                    )
                ]
            )
        ),
        links: links,
        assets: assets
    }::JSON
FROM '../data/collections.jsonl';
