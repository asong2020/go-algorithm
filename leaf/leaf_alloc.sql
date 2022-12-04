create table leaf.leaf_alloc
(
	id int auto_increment
		primary key,
	biz_tag varchar(128) default '' not null,
	max_id bigint default 1 not null,
	step int not null,
	description varchar(256) null,
	update_time bigint default 0 not null,
	constraint biz_tag
		unique (biz_tag)
);

