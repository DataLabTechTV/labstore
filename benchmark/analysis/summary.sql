CREATE OR REPLACE TABLE summary_stats AS
WITH bandwidth AS (
    SELECT
        median(up_bytes_per_sec) AS median_up_bytes_per_sec,
        median(down_bytes_per_sec) AS median_down_bytes_per_sec
    FROM iperf3_stats
),
stats AS (
    SELECT
        op,
        store,
        median(obj_per_sec) AS median_obj_per_sec,
        median(bytes_per_sec / 1024 / 1024) AS median_mebibytes_per_sec,
        median(
            CASE
                WHEN op = 'PUT' THEN bytes_per_sec / median_up_bytes_per_sec
                WHEN op = 'GET' THEN bytes_per_sec / median_down_bytes_per_sec
                ELSE 0
            END
        ) AS bandwidth_pct
    FROM warp_stats, bandwidth
    GROUP BY store, op
    ORDER BY op, median_obj_per_sec DESC
)
SELECT row_number() OVER (PARTITION BY op ORDER BY median_obj_per_sec DESC) AS rank, *
FROM stats;

SELECT * FROM summary_stats WHERE op = 'DELETE';
SELECT * FROM summary_stats WHERE op = 'GET';
SELECT * FROM summary_stats WHERE op = 'PUT';
SELECT * FROM summary_stats WHERE op = 'STAT';

COPY (SELECT * FROM summary_stats WHERE op = 'DELETE') TO 'analysis/output/op_stats-delete.dat' (DELIMITER '\t');
COPY (SELECT * FROM summary_stats WHERE op = 'GET') TO 'analysis/output/op_stats-get.dat' (DELIMITER '\t');
COPY (SELECT * FROM summary_stats WHERE op = 'PUT') TO 'analysis/output/op_stats-put.dat' (DELIMITER '\t');
COPY (SELECT * FROM summary_stats WHERE op = 'STAT') TO 'analysis/output/op_stats-stat.dat' (DELIMITER '\t');
