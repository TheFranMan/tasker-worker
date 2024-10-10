TRUNCATE TABLE jobs;
TRUNCATE TABLE requests;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('2b482d15-6c02-4e7f-bae3-0a8fe1dfb301', 'auth-token-valid-1', 'delete', '{\"id\": 2}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 0, '2024-10-09 09:13:59', NULL),
	('482a2d88-d38a-4509-ac94-beadff53c053', 'auth-token-valid-1', 'delete', '{\"id\": 1}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 0, '2024-10-09 09:13:56', NULL),
	('7fbef510-e37d-4884-97e2-c31fac6a89ae', 'auth-token-valid-1', 'delete', '{\"id\": 100}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 3, '2024-10-09 09:11:10', '2024-10-09 09:11:10');
INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', '7fbef510-e37d-4884-97e2-c31fac6a89ae', 0, 'test error', 3, '2024-10-09 09:12:00', '2024-10-09 09:13:00');