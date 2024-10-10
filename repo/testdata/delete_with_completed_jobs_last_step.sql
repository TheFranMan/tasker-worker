TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', 'f51d7890-8d28-4f9a-9803-6f2645c7d6c7', 0, 'test error', 3, '2024-10-09 16:01:00', '2024-10-09 16:02:00'),
	(2, 'service1GetUser', '80498d81-8de4-41fb-b1a5-53180cd56d73', 0, NULL, 2, '2024-10-10 11:29:00', '2024-10-10 11:30:00'),
	(3, 'service1DeleteUser', '80498d81-8de4-41fb-b1a5-53180cd56d73', 1, NULL, 2, '2024-10-10 11:31:00', '2024-10-10 11:32:00'),
	(4, 'service2DeleteUser', '80498d81-8de4-41fb-b1a5-53180cd56d73', 1, NULL, 2, '2024-10-10 11:31:00', '2024-10-10 11:32:00'),
	(5, 'service3DeleteUser', '80498d81-8de4-41fb-b1a5-53180cd56d73', 1, NULL, 2, '2024-10-10 11:31:00', '2024-10-10 11:32:00');
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('80498d81-8de4-41fb-b1a5-53180cd56d73', 'auth-token-valid-1', 'delete', '{\"id\": 1}', '{\"email\": \"example_1@example.com\"}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 1, 1, '2024-10-10 11:28:12', NULL),
	('f51d7890-8d28-4f9a-9803-6f2645c7d6c7', 'auth-token-valid-1', 'delete', '{\"id\": 100}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 3, '2024-10-09 16:00:14', '2024-10-09 16:03:00');