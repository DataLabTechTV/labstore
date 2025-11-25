WITH data AS (
    SELECT
        'DELETE' AS op,
        "DELETE".*
    FROM
        read_json(
            'warp/output/docker-apps_8333/*.json.zst',
            records=true
        ) j,
        UNNEST(j.by_op_type.DELETE.throughput.segmented.segments) "DELETE",
        UNNEST(j.by_op_type.GET.throughput.segmented.segments) "GET"
)
SELECT *
FROM data;
