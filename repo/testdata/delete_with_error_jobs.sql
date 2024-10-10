TRUNCATE TABLE requests;
TRUNCATE TABLE jobs;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('7fbef510-e37d-4884-97e2-c31fac6a89ae', 'auth-token-valid-1', 'delete', '{\"id\": 100}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 3, '2024-10-09 09:11:10', '2024-10-09 09:11:10'),
	('89858b95-21bd-47e3-a03e-9069a7440188', 'auth-token-valid-1', 'delete', '{\"id\": 1}', NULL, '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service1DeleteUser\", \"service2DeleteUser\", \"service3DeleteUser\"], \"name\": \"delete_user_accounts\"}]', 0, 1, '2024-10-09 13:37:28', NULL);

INSERT INTO `jobs` (`id`, `name`, `token`, `step`, `error`, `status`, `created`, `completed`)
VALUES
	(1, 'service1GetUser', '7fbef510-e37d-4884-97e2-c31fac6a89ae', 0, 'test error', 3, '2024-10-09 09:12:00', '2024-10-09 09:13:00'),
	(2, 'service1GetUser', '89858b95-21bd-47e3-a03e-9069a7440188', 0, 'test error', 4, '2024-10-09 13:38:00', '2024-10-09 13:39:00');

