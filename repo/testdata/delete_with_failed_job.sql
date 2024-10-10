TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('c639b525-1ab1-44e6-bde3-96238cf13f2f', 'auth-token-valid-1', 'delete', '{\"id\": 1}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 1, '2024-10-09 15:49:39', NULL),
	('f0c8c981-b518-4215-9fa3-804fd0dc2ba1', 'auth-token-valid-1', 'delete', '{\"id\": 100}', '{\"email\": \"example_100@example.com\"}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 1, 2, '2024-10-09 15:43:35', '2024-10-09 15:48:00');
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', 'f0c8c981-b518-4215-9fa3-804fd0dc2ba1', 0, NULL, 2, '2024-10-09 15:44:00', '2024-10-09 15:45:00'),
	(2, 'service1DeleteUser', 'f0c8c981-b518-4215-9fa3-804fd0dc2ba1', 1, NULL, 2, '2024-10-09 15:46:00', '2024-10-09 15:47:00'),
	(3, 'service2DeleteUser', 'f0c8c981-b518-4215-9fa3-804fd0dc2ba1', 1, NULL, 2, '2024-10-09 15:46:00', '2024-10-09 15:47:00'),
	(4, 'service3DeleteUser', 'f0c8c981-b518-4215-9fa3-804fd0dc2ba1', 1, NULL, 2, '2024-10-09 15:46:00', '2024-10-09 15:47:00'),
	(5, 'service1GetUser', 'c639b525-1ab1-44e6-bde3-96238cf13f2f', 0, 'test error', 3, '2024-10-09 15:50:00', '2024-10-09 15:51:00');