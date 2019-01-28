
BEGIN;

ALTER TABLE `application` DROP COLUMN `klinks`;

DROP TABLE `klink`;

COMMIT;