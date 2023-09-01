DROP TABLE IF EXISTS user_segments;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS segments;

update user_segments us
set deleted_at = NULL
where segment_id in (select id from segments where slug = 'test')
  and us.user_id = 7
  and deleted_at is not null;

update user_segments us
set deleted_at = NOW()
where segment_id in (select id from segments where slug = 'test')
  and us.user_id = 6
  and deleted_at is null;

insert into user_segments (user_id, segment_id)
values (6, 1)
on conflict (user_id, segment_id) do update set deleted_at = NULL
where user_segments.segment_id in (select id from segments where slug = 'test')
  and user_segments.user_id = 6
  and user_segments.deleted_at is not null;

