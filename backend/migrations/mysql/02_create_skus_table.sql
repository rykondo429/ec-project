-- SKUテーブル作成マイグレーション
-- 商品バリエーション（サイズ、色など）ごとの在庫管理

CREATE TABLE IF NOT EXISTS skus (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL,
    sku_code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    attributes JSON,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    stock INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_product_id (product_id),
    INDEX idx_sku_code (sku_code),
    INDEX idx_is_active (is_active),
    INDEX idx_product_sku (product_id, is_active),
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 在庫履歴テーブル（将来的な拡張用）
CREATE TABLE IF NOT EXISTS sku_stock_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sku_id VARCHAR(36) NOT NULL,
    quantity_change INT NOT NULL,
    stock_before INT NOT NULL,
    stock_after INT NOT NULL,
    reason VARCHAR(255),
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_sku_id (sku_id),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (sku_id) REFERENCES skus (id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- サンプルデータ: Nike Air Max 90のSKU
INSERT INTO
    skus (
        id,
        product_id,
        sku_code,
        name,
        attributes,
        price,
        stock,
        is_active
    )
VALUES (
        'sku-001',
        'prod-001',
        'NIKE-AM90-RED-S',
        'Nike Air Max 90 - Red / S',
        '{"color": "Red", "size": "S"}',
        12000,
        15,
        TRUE
    ),
    (
        'sku-002',
        'prod-001',
        'NIKE-AM90-RED-M',
        'Nike Air Max 90 - Red / M',
        '{"color": "Red", "size": "M"}',
        12000,
        20,
        TRUE
    ),
    (
        'sku-003',
        'prod-001',
        'NIKE-AM90-RED-L',
        'Nike Air Max 90 - Red / L',
        '{"color": "Red", "size": "L"}',
        12000,
        15,
        TRUE
    ),
    (
        'sku-004',
        'prod-001',
        'NIKE-AM90-BLK-S',
        'Nike Air Max 90 - Black / S',
        '{"color": "Black", "size": "S"}',
        12000,
        10,
        TRUE
    ),
    (
        'sku-005',
        'prod-001',
        'NIKE-AM90-BLK-M',
        'Nike Air Max 90 - Black / M',
        '{"color": "Black", "size": "M"}',
        12000,
        25,
        TRUE
    ),
    (
        'sku-006',
        'prod-001',
        'NIKE-AM90-BLK-L',
        'Nike Air Max 90 - Black / L',
        '{"color": "Black", "size": "L"}',
        12000,
        20,
        TRUE
    );

-- サンプルデータ: Adidas Ultraboost 22のSKU
INSERT INTO
    skus (
        id,
        product_id,
        sku_code,
        name,
        attributes,
        price,
        stock,
        is_active
    )
VALUES (
        'sku-007',
        'prod-002',
        'ADIDAS-UB22-WHT-S',
        'Adidas Ultraboost 22 - White / S',
        '{"color": "White", "size": "S"}',
        18000,
        8,
        TRUE
    ),
    (
        'sku-008',
        'prod-002',
        'ADIDAS-UB22-WHT-M',
        'Adidas Ultraboost 22 - White / M',
        '{"color": "White", "size": "M"}',
        18000,
        12,
        TRUE
    ),
    (
        'sku-009',
        'prod-002',
        'ADIDAS-UB22-WHT-L',
        'Adidas Ultraboost 22 - White / L',
        '{"color": "White", "size": "L"}',
        18000,
        10,
        TRUE
    );

-- サンプルデータ: Uniqlo Heattech（サイズバリエーション）
INSERT INTO
    skus (
        id,
        product_id,
        sku_code,
        name,
        attributes,
        price,
        stock,
        is_active
    )
VALUES (
        'sku-010',
        'prod-010',
        'UNIQLO-HT-BLK-S',
        'ユニクロ ヒートテック - Black / S',
        '{"color": "Black", "size": "S"}',
        1990,
        40,
        TRUE
    ),
    (
        'sku-011',
        'prod-010',
        'UNIQLO-HT-BLK-M',
        'ユニクロ ヒートテック - Black / M',
        '{"color": "Black", "size": "M"}',
        1990,
        50,
        TRUE
    ),
    (
        'sku-012',
        'prod-010',
        'UNIQLO-HT-BLK-L',
        'ユニクロ ヒートテック - Black / L',
        '{"color": "Black", "size": "L"}',
        1990,
        35,
        TRUE
    ),
    (
        'sku-013',
        'prod-010',
        'UNIQLO-HT-WHT-S',
        'ユニクロ ヒートテック - White / S',
        '{"color": "White", "size": "S"}',
        1990,
        15,
        TRUE
    ),
    (
        'sku-014',
        'prod-010',
        'UNIQLO-HT-WHT-M',
        'ユニクロ ヒートテック - White / M',
        '{"color": "White", "size": "M"}',
        1990,
        25,
        TRUE
    );