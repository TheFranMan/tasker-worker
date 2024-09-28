TRUNCATE TABLE jobs;
TRUNCATE TABLE requests;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('test-token-1', 'auth-token-1', 'Delete', '{\"id\": 1}', '{}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 0, '2024-09-28 14:14:27', NULL),
	('test-token-2', 'auth-token-2', 'Delete', '{\"id\": 2}', '{}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 0, '2024-09-28 14:16:54', NULL),
	('test-token-3', 'auth-token-3', 'Delete', '{\"id\": 4}', '{}', '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 6, '2024-09-28 14:17:01', NULL);
