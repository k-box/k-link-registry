
BEGIN;

UPDATE `registrant` SET last_login = 0 WHERE last_login IS NULL;

ALTER TABLE `registrant` MODIFY `last_login` int(11) NOT NULL DEFAULT 0;

COMMIT;