-- =============================================
-- Load warp statistics
-- =============================================

CREATE OR REPLACE MACRO warp_load(folder) AS TABLE
SELECT
    op.key AS op,
    (s.value->>'start')::TIMESTAMP AS start,
    (s.value->>'bytes_per_sec')::INTEGER AS bytes_per_sec,
    (s.value->>'obj_per_sec')::FLOAT AS obj_per_sec
FROM
    read_json(folder || '/*.json.zst', records=true) j,
    json_each(j.by_op_type) op,
    json_each(op.value.throughput.segmented.segments) s;

CREATE OR REPLACE TABLE warp_stats AS
SELECT 'LabStore' AS store, labstore.*,
FROM warp_load('warp/output/labstore') labstore
UNION
SELECT 'MinIO' AS store, minio.*
FROM warp_load('warp/output/minio') minio
UNION
SELECT 'Garage' AS store, garage.*
FROM warp_load('warp/output/garage') garage
UNION
SELECT 'SeaweedFS' AS store, seaweedfs.*
FROM warp_load('warp/output/seaweedfs') seaweedfs
UNION
SELECT 'RustFS' AS store, rustfs.*
FROM warp_load('warp/output/rustfs') rustfs
ORDER BY store, op, start;

-- =============================================
-- Load iperf3 statistics
-- =============================================

CREATE OR REPLACE TABLE iperf3_stats AS
SELECT
    make_timestamp_ms(j.start.timestamp.timesecs * 1000) AS start,
    j.end.sum_received.bits_per_second / 8 AS down_bytes_per_sec,
    j.end.sum_sent.bits_per_second / 8 AS up_bytes_per_sec
FROM
    read_json('iperf3/output/*.json.gz', records=true) j;
