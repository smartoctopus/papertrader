-- migrate:up
-- Instruments metadata
CREATE TABLE instruments (
    instrument_id INTEGER PRIMARY KEY,
    symbol TEXT NOT NULL UNIQUE
);
CREATE UNIQUE INDEX idx_instruments_id_symbol ON instruments (instrument_id, symbol);

-- Tick data
CREATE TABLE ticks (
    tick_id INTEGER PRIMARY KEY,
    instrument_id INTEGER NOT NULL,

    time INTEGER NOT NULL, -- Unix milliseconds timestamp
    price DECIMAL NOT NULL,
    volume DECIMAL NOT NULL,

    FOREIGN KEY (instrument_id) REFERENCES instruments(instrument_id)
);
CREATE UNIQUE INDEX idx_ticks_instrument_time ON ticks (instrument_id, time);

-- OHLCV Data
CREATE TABLE ohlcv (
    ohlcv_id INTEGER PRIMARY KEY,
    instrument_id INTEGER NOT NULL,

    time INTEGER NOT NULL, -- Unix timestamp
    open DECIMAL NOT NULL,
    high DECIMAL NOT NULL,
    low DECIMAL NOT NULL,
    close DECIMAL NOT NULL,
    volume DECIMAL NOT NULL,

    FOREIGN KEY (instrument_id) REFERENCES instruments(instrument_id)
);
CREATE UNIQUE INDEX idx_ohlcv_instrument_interval_time ON ohlcv (instrument_id, time);

-- migrate:down
DROP TABLE instruments
DROP TABLE ticks
DROP TABLE ohlcv
