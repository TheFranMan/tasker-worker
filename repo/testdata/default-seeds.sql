TRUNCATE TABLE jobs;
TRUNCATE TABLE requests;
INSERT INTO `requests` (`token`, `request_token`, `action`, `params`, `extras`, `steps`, `step`, `status`, `created`, `completed`)
VALUES
	('lala', 'hoho', 'delete', '{\"id\": 1}', "{}", '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service2DeleteUser\", \"service1DeleteUser\"], \"name\": \"service2_delete_user\"}]', 0, 0, '2024-09-09 14:23:21', NULL),
	('lala2', 'hoho2', 'delete2', '{\"id\": 2}', "{}", '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service2DeleteUser\", \"service1DeleteUser\"], \"name\": \"service2_delete_user\"}]', 0, 0, '2024-09-09 14:23:26', NULL),
	('lala3', 'hoho3', 'delete3', '{\"id\": 3}', "{}", '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service2DeleteUser\", \"service1DeleteUser\"], \"name\": \"service2_delete_user\"}]', 0, 1, '2024-09-09 14:23:30', NULL),
	('lala4', 'hoho4', 'delete4', '{\"id\": 4}', "{}", '[{\"jobs\": [\"service1GetUser\"], \"name\": \"service1_retrieve_user\"}, {\"jobs\": [\"service2DeleteUser\", \"service1DeleteUser\"], \"name\": \"service2_delete_user\"}]', 0, 6, '2024-09-09 14:23:34', NULL);