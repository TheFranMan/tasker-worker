CREATE TABLE `requests` (
  `token` varchar(255) NOT NULL DEFAULT '',
  `request_token` varchar(20) NOT NULL DEFAULT '',
  `action` varchar(255) NOT NULL DEFAULT '',
  `params` json NOT NULL,
  `extras` json DEFAULT NULL,
  `steps` json NOT NULL,
  `step` int NOT NULL DEFAULT '0',
  `status` int NOT NULL DEFAULT '0',
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed` datetime DEFAULT NULL,
  PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;