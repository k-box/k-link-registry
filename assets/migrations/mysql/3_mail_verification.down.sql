
BEGIN;

--
-- Table structure for table `email_verification`
--
CREATE TABLE IF NOT EXISTS `email_verification` (
  `email` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `registrant_id` bigint(20) NOT NULL,
  `token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`email`)
);

--
-- Table structure for table `password_change_verification`
--
CREATE TABLE IF NOT EXISTS `password_change_verification` (
  `registrant_id` bigint(20) NOT NULL,
  `token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `timestamp` int(11) NOT NULL,
  PRIMARY KEY (`registrant_id`)
);

ALTER TABLE `registrant`
    DROP COLUMN `confirmed`,
    DROP COLUMN `confirm_selector`,
    DROP COLUMN `confirm_verifier`,
    DROP COLUMN `recover_selector`,
    DROP COLUMN `recover_verifier`,
    DROP COLUMN `recover_expiry`;

COMMIT;