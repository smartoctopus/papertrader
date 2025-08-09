-- name: CreateInstrument :one
INSERT OR IGNORE INTO instruments (
    symbol
) VALUES (?)
RETURNING *;

-- name: GetInstrument :one
SELECT * FROM instruments
WHERE symbol = ? LIMIT 1;

-- name: InsertOHLCV :exec
INSERT OR IGNORE INTO ohlcv (
    instrument_id,
    time,
    open,
    high,
    low,
    close,
    volume
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: InsertTick :exec
INSERT OR IGNORE INTO ticks (
    instrument_id,
    time,
    price,
    volume
) VALUES (?, ?, ?, ?);

-- name: GetOHLCV :many
SELECT o.time, o.open, o.high, o.low, o.close, o.volume
FROM ohlcv o
JOIN instruments i ON o.instrument_id = i.instrument_id
WHERE o.time BETWEEN :start AND :end
    AND i.symbol = :instrument
ORDER BY o.time;

-- name: GetTicks :many
SELECT t.time, t.price, t.volume
FROM ticks t
JOIN instruments i ON t.instrument_id = i.instrument_id
WHERE t.time BETWEEN :start AND :end
    AND i.symbol = :instrument
ORDER BY t.time;
