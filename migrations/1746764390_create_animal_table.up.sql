-- Migration Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE IF NOT EXISTS `animals` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `species` varchar(255) NOT NULL,
  `age` bigint,
  `description` varchar(255),
  `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_animal_name` (`name`),
  KEY `idx_animal_species` (`species`),
  KEY `idx_animal_age` (`age`),
  KEY `idx_animal_created_at` (`created_at`),
  KEY `idx_animal_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;