begin;

insert into users (guid, username, telegram_id)
values ('00000000-0000-0000-0000-000000000001', 'username1', '00000001');

insert into spending_categories (guid, user_guid, category, description)
values ('00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000001', 'category1', 'bla bla bla'),
       ('00000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000001', 'category2', 'bla bla bla'),
       ('00000000-0000-0000-0000-000000000031', '00000000-0000-0000-0000-000000000001', 'category3', 'bla bla bla'),
       ('00000000-0000-0000-0000-000000000041', '00000000-0000-0000-0000-000000000001', 'category4', 'bla bla bla');

commit;