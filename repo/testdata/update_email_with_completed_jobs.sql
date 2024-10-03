TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('test-token-1', 'auth-token-1', 'update_email', '{\"id\": 1}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 1, '2024-09-30 15:29:03', NULL),
	('test-token-2', 'auth-token-2', 'update_email', '{\"id\": 2}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 1, '2024-09-30 15:29:03', NULL),
	('test-token-3', 'auth-token-3', 'update_email', '{\"id\": 4}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 6, '2024-09-28 14:34:30', NULL);
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1UpdateUser', 'test-token-1', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00'),
	(2, 'service2UpdateUser', 'test-token-1', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00'),
	(3, 'service3UpdateUser', 'test-token-1', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00'),
	(4, 'service1UpdateUser', 'test-token-2', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00'),
	(5, 'service2UpdateUser', 'test-token-2', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00'),
	(6, 'service3UpdateUser', 'test-token-2', 0, NULL, 2, '2024-09-30 15:29:03', '2024-10-03 21:07:00');