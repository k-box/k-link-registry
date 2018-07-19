-- This migration initializes the legacy version of the k-registry, so that
-- more modern revisions can be based on top of it.

-- CAVEAT: All tables are created using the IF NOT EXISTS directive, as
-- the tables might have already been created by the legacy k-registry,
-- and no migration table was created yet.
-- This ensures that the migration will apply, and the state will be
-- consistent afterwards

BEGIN;

--
-- Table structure for table `registrant`
--
CREATE TABLE IF NOT EXISTS `registrant` (
  `registrant_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `password` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `salt` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL,
  `last_login` int(11) DEFAULT NULL,
  PRIMARY KEY (`registrant_id`),
  UNIQUE KEY (`email`)
);

--
-- Table structure for table `application`
--
CREATE TABLE IF NOT EXISTS `application` (
  `application_id` int(11) NOT NULL AUTO_INCREMENT,
  `registrant_id` bigint(20) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `app_domain` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  `auth_token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `permissions` longtext COLLATE utf8mb4_unicode_ci COMMENT '(DC2Type:simple_array)',
  `status` tinyint(1) NOT NULL,
  PRIMARY KEY (`application_id`),
  UNIQUE KEY (`app_domain`),
  KEY (`registrant_id`),
  CONSTRAINT FOREIGN KEY (`registrant_id`) REFERENCES `registrant` (`registrant_id`)
);

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

--
-- Table structure for table `permission`
--
CREATE TABLE IF NOT EXISTS `permission` (
  `name` varchar(150) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`name`),
  UNIQUE KEY (`name`)
);

COMMIT;
