
BEGIN;

ALTER TABLE `registrant` MODIFY `last_login` int(11) NULL;

UPDATE `registrant` SET last_login = NULL WHERE last_login=0;

COMMIT;