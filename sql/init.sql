alter table entry add link text;

CREATE TABLE `cached_content` (
`keyword` varchar(255) DEFAULT NULL,
`content` mediumtext,
KEY `keyword_index` (`keyword`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin