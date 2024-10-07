TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('test-token-1', 'auth-token-1', 'Delete', '{\"id\": 1}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 1, '2024-10-01 13:14:33', NULL),
	('test-token-2', 'auth-token-2', 'Delete', '{\"id\": 2}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 1, '2024-10-01 13:14:33', NULL),
	('test-token-3', 'auth-token-3', 'Delete', '{\"id\": 4}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 6, '2024-09-28 14:17:01', NULL);
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', 'test-token-1', 0, NULL, 0, '2024-10-01 13:14:33', NULL),
	(2, 'service1GetUser', 'test-token-2', 0, NULL, 0, '2024-10-01 13:14:33', NULL);