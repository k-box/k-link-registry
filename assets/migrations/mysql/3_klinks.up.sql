-- This migration add support for handling multiple K-Links, 
-- as per the merge networks concept

BEGIN;

--
-- Table structure for table `klink`
--
CREATE TABLE IF NOT EXISTS `klink` (
  `klink_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `identifier` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, -- the public K-Link identifier
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `website` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` longtext COLLATE utf8mb4_unicode_ci DEFAULT NULL, -- a description of the network
  `manager_id` bigint(20) DEFAULT NULL,
  `active` tinyint(1) NOT NULL, -- indicate if the K-Link can be selected for applications or should only be visible to managers
  PRIMARY KEY (`klink_id`),
  UNIQUE KEY (`identifier`),
  KEY (`manager_id`),
  CONSTRAINT FOREIGN KEY (`manager_id`) REFERENCES `registrant` (`registrant_id`)
);

--
-- Add the klinks column to the applications table to store the klinks to which the application can publish
-- like it was done for the permissions
--
ALTER TABLE `application` ADD COLUMN `klinks` longtext COLLATE utf8mb4_unicode_ci COMMENT '(DC2Type:simple_array)';

COMMIT;
