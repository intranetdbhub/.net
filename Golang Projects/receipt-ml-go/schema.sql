-- schema.sql
CREATE DATABASE IF NOT EXISTS receiptsdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE receiptsdb;

CREATE TABLE IF NOT EXISTS receipts (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  merchant VARCHAR(255) NULL,
  purchase_date DATE NULL,
  subtotal DECIMAL(10,2) NULL,
  tax DECIMAL(10,2) NULL,
  fees DECIMAL(10,2) NULL,
  total DECIMAL(10,2) NULL,
  currency VARCHAR(10) NULL,
  not_visible JSON NULL,
  ocr_text LONGTEXT NOT NULL,
  original_filename VARCHAR(255),
  storage_path VARCHAR(512),
  status ENUM('parsed','needs_review') NOT NULL DEFAULT 'needs_review',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS receipt_items (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  receipt_id BIGINT NOT NULL,
  label VARCHAR(255) NOT NULL,
  quantity DECIMAL(10,2) DEFAULT 1.00,
  unit_price DECIMAL(10,2) NULL,
  line_total DECIMAL(10,2) NULL,
  is_fee BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE
) ENGINE=InnoDB;
