-- Migration Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS `flowers` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `species` varchar(255) NOT NULL,
    `color` varchar(50) NOT NULL,
    `description` varchar(1000),
    `seasonal` boolean DEFAULT false,
    `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    KEY `idx_flower_name` (`name`),
    KEY `idx_flower_species` (`species`),
    KEY `idx_flower_color` (`color`),
    KEY `idx_flower_created_at` (`created_at`),
    KEY `idx_flower_updated_at` (`updated_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;