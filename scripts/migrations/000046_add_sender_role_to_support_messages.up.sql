ALTER TABLE `support_messages`
ADD COLUMN `sender_role` VARCHAR(50) NULL AFTER `sender_name`;
