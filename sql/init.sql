-- データベースの文字セットを明示的に設定
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;
SET CHARACTER SET utf8mb4;

-- Create items table for managing valuable items and collections
CREATE TABLE IF NOT EXISTS items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT 'Item name',
    category VARCHAR(50) NOT NULL COMMENT 'Item category: 時計, バッグ, ジュエリー, 靴, その他',
    brand VARCHAR(100) NOT NULL COMMENT 'Brand name',
    purchase_price INT NOT NULL DEFAULT 0 COMMENT 'Purchase price in yen',
    purchase_date DATE NOT NULL COMMENT 'Purchase date in YYYY-MM-DD format',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp',
    
    INDEX idx_category (category),
    INDEX idx_brand (brand),
    INDEX idx_purchase_date (purchase_date),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Table for managing valuable items and collections';

-- Insert sample data for testing
INSERT INTO items (name, category, brand, purchase_price, purchase_date) VALUES
('ロレックス デイトナ', '時計', 'ROLEX', 1500000, '2023-01-15'),
('エルメス バーキン', 'バッグ', 'HERMÈS', 2000000, '2023-02-20'),
('ティファニー ネックレス', 'ジュエリー', 'Tiffany & Co.', 300000, '2023-03-10'),
('ルブタン パンプス', '靴', 'Christian Louboutin', 150000, '2023-04-05'),
('アップルウォッチ', 'その他', 'Apple', 50000, '2023-05-12');