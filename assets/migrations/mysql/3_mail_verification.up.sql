
BEGIN;

ALTER TABLE `registrant`
    ADD COLUMN `confirmed` BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN `confirm_selector` varchar(255) NOT NULL DEFAULT "",
    ADD COLUMN `confirm_verifier` varchar(255) NOT NULL DEFAULT "",
    ADD COLUMN `recover_selector` varchar(255) NOT NULL DEFAULT "",
    ADD COLUMN `recover_verifier` varchar(255) NOT NULL DEFAULT "",
    ADD COLUMN `recover_expiry` bigint NOT NULL DEFAULT 0;

DROP TABLE `email_verification`;

DROP TABLE `password_change_verification`;

COMMIT;