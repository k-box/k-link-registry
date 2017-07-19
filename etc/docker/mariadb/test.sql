USE `kregistry`;
INSERT IGNORE INTO `registrant`(`registrant_id`, `password`, `salt`, `email`, `name`, `role`, `status`, `last_login`) VALUES(1, NULL, NULL, '', '', '', 0, 0);
INSERT IGNORE INTO `application`(`application_id`, `registrant_id`, `name`, `app_domain`, `auth_token`, `permissions`, `status`) VALUES(1, 1, 'test', 'https://kregistry.core', '0245754939da6d7576c52f72ad7ac7a76932a0fc62f845f3d5ed80b042df2442858b88a8d51cc73d9c030c3f5577c2b0c91f56bf5001312734b1d7ee25921a6d', 'data_view,data_add,data_edit', 1);
