--
-- File generated with SQLiteStudio v3.2.1 on Mon May 13 22:58:48 2019
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: item_amount
CREATE TABLE item_amount (sku VARCHAR (30) NOT NULL PRIMARY KEY, name VARCHAR (100) NOT NULL, quantity INTEGER DEFAULT (0));

-- Table: item_in
CREATE TABLE item_in (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, time DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), sku VARCHAR (30) REFERENCES item_amount (sku) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL, amount_orders INTEGER DEFAULT (0), amount_received INTEGER DEFAULT (0), purchase_price INTEGER DEFAULT (0), receipt_number VARCHAR (30), notes STRING);

-- Table: item_out
CREATE TABLE item_out (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, time DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), sku VARCHAR (30) REFERENCES item_amount (sku) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL, amount_out INTEGER DEFAULT (0), selling_price INTEGER DEFAULT (0), order_id VARCHAR (30), notes TEXT);

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
