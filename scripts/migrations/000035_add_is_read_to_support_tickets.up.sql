CREATE TABLE IF NOT EXISTS `support_ticket_reads` (
  `ticket_id` CHAR(36) NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `read_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ticket_id`, `user_id`),
  CONSTRAINT `fk_support_ticket_reads_ticket` FOREIGN KEY (`ticket_id`)
    REFERENCES `support_tickets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_support_ticket_reads_user` FOREIGN KEY (`user_id`)
    REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
