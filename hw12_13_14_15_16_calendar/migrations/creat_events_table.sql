
create schema if not exists calendar;
    drop table if exists events;
     create table events (
         id serial primary key,
         title varchar(50) not null,
         date_and_time_of_the_event timestamp default current_timestamp,
         duration_of_the_event timestamp not null,
         description_event varchar(255),
         id_user bigint not null,
         time_until_event bigint
     );
    create index idx_events on events(id);