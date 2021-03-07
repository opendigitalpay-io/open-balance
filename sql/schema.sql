CREATE TABLE IF NOT EXISTS `users`
(
    `id`          bigint(11)   NOT NULL,
    `email`       varchar(255) NOT NULL,
    `phone`       varchar(255) NOT NULL,
    `external_id` varchar(255) NOT NULL,
    `metadata`    json                  DEFAULT NULL,
    `created_at`  bigint(11)   NOT NULL,
    `updated_at`  bigint(11)   NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
)
    DEFAULT CHARSET = utf8mb4
;

CREATE TABLE IF NOT EXISTS `root_accounts`
(
    `id`         bigint(11)   NOT NULL,
    `user_id`    bigint(11)   NOT NULL,
    `type`       varchar(255) NOT NULL,
    `state`      varchar(255) NOT NULL,
    `metadata`   text                  DEFAULT NULL,
    `created_at` bigint(11)   NOT NULL,
    `updated_at` bigint(11)   NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_id` (`user_id`)
)
    DEFAULT CHARSET = utf8mb4
;

CREATE TABLE IF NOT EXISTS `balance_accounts`
(
    `id`              bigint(11)   NOT NULL,
    `root_account_id` bigint(11)   NOT NULL,
    `type`            varchar(255) NOT NULL,
    `state`           varchar(255) NOT NULL,
    `visible`         boolean      NOT NULL DEFAULT TRUE,
    `lockable`        boolean      NOT NULL DEFAULT TRUE,
    `balance`         bigint(11)   NOT NULL default '0',
    `currency`        varchar(255) NOT NULL,
    `version`         int(11)      NOT NULL DEFAULT '1',
    `metadata`        text                  DEFAULT NULL,
    `created_at`      bigint(11)   NOT NULL,
    `updated_at`      bigint(11)   NOT NULL,
    PRIMARY KEY (`id`),
    KEY `index_balance_accounts_on_root_account_id` (`root_account_id`) USING BTREE
)
    DEFAULT CHARSET = utf8mb4
;

CREATE TABLE `transactions`
(
    `id`               bigint(11)   NOT NULL,
    `parent_id`        bigint(11),
    `src_account_id`   bigint(11)   NOT NULL,
    `dst_account_id`   bigint(11)   NOT NULL,
    `src_user_id`      bigint(11)   NOT NULL,
    `dst_user_id`      bigint(11)   NOT NULL,
    `amount`           bigint(11)   NOT NULL,
    `currency`         varchar(255) NOT NULL,
    `src_balance`      bigint(11)   NOT NULL,
    `dst_balance`      bigint(11)   NOT NULL,
    `src_account_type` varchar(255) NOT NULL,
    `dst_account_type` varchar(255) NOT NULL,
    `reversible`       boolean      NOT NULL DEFAULT TRUE,
    `metadata`         text                  DEFAULT NULL,
    `created_at`       bigint(11)   NOT NULL,
    PRIMARY KEY (`id`),
    KEY `index_transactions_on_src_account_id` (`src_account_id`) USING BTREE,
    KEY `index_transactions_on_dst_account_id` (`dst_account_id`) USING BTREE,
    KEY `index_transactions_on_src_user_id` (`src_user_id`) USING BTREE,
    KEY `index_transactions_on_dst_user_id` (`dst_user_id`) USING BTREE
)
    DEFAULT CHARSET = utf8mb4
;

CREATE TABLE `idempotency`
(
    `id`           bigint(11)   NOT NULL,
    `idem_id`      varchar(255) NOT NULL,
    `is_completed` boolean      NOT NULL DEFAULT FALSE,
    `response`     text         DEFAULT NULL,
    `created_at`   bigint(11)   NOT NULL,
    `updated_at`   bigint(11)   NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idem_id` (`idem_id`)
)
    DEFAULT CHARSET = utf8mb4
;

INSERT INTO `users` VALUES
('1', 'open-balance-system@gmail.com', '6471111111', 'open-balance-system', '{\"rootAccess\": true}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('2', 'open-balance@gmail.com', '6471111111', 'open-balance', '{\"rootAccess\": true}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('3', 'open-balance-biz@gmail.com', '6471111111', 'open-balance-biz', '{\"rootAccess\": true}', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO `root_accounts` VALUES
('1', '1', 'SYSTEM', 'ACTIVE', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('2', '2', 'PERSONAL', 'ACTIVE', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('3', '3', 'PERSONAL', 'ACTIVE', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO `balance_accounts` VALUES
('1', '1', 'SOURCE', 'ACTIVE', FALSE, TRUE, '10000000', 'CAD', '1', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('2', '2', 'CHEQUE', 'ACTIVE', TRUE, TRUE, '5000', 'CAD', '1', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('3', '2', 'PAYMENT', 'ACTIVE', FALSE, TRUE, '0', 'CAD', '1', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('4', '3', 'PAYABLE', 'ACTIVE', FALSE, TRUE, '0', 'CAD', '1', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('5', '2', 'INCOMING', 'ACTIVE', FALSE, TRUE, '0', 'CAD', '1', null, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
