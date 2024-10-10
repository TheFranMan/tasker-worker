TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('039a4e90-107b-4d7f-97f7-e1ad84316119', 'auth-token-valid-1', 'delete', '{\"id\": 1}', '{\"email\": \"example_1@example.com\"}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 1, '2024-10-09 16:04:53', NULL),
	('f51d7890-8d28-4f9a-9803-6f2645c7d6c7', 'auth-token-valid-1', 'delete', '{\"id\": 100}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 3, '2024-10-09 16:00:14', '2024-10-09 16:03:00');
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', 'f51d7890-8d28-4f9a-9803-6f2645c7d6c7', 0, 'test error', 3, '2024-10-09 16:01:00', '2024-10-09 16:02:00'),
	(2, 'service1GetUser', '039a4e90-107b-4d7f-97f7-e1ad84316119', 0, NULL, 2, '2024-10-09 16:05:00', '2024-10-09 16:06:00');
