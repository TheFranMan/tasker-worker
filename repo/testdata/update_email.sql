TRUNCATE TABLE jobs;
TRUNCATE TABLE requests;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('test-token-1', 'auth-token-1', 'update_email', '{\"id\": 1}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 0, '2024-09-28 14:34:21', NULL),
	('test-token-2', 'auth-token-2', 'update_email', '{\"id\": 2}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 0, '2024-09-28 14:34:25', NULL),
	('test-token-3', 'auth-token-3', 'update_email', '{\"id\": 4}', '{}', '[{\"jobs\": [\"service1UpdateUser\", \"service2UpdateUser\", \"service3UpdateUser\"], \"name\": \"update_emails\"}]', 0, 6, '2024-09-28 14:34:30', NULL);
