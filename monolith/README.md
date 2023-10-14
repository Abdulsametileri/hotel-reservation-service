### Creating Reservation Table
```sql
create table reservation
(
    reservation_id uuid
        constraint reservation_pk
            primary key,
    hotel_id       integer,
    room_type_id   integer,
    start_date     date,
    end_date       date,
    status         text
);
```

Status can be in one of these states: pending, paid, refunded, canceled, rejected.


### Creating Room Type Inventory Table
```sql
create table room_type_inventory
(
    hotel_id        integer,
    room_type_id    integer,
    date            date,
    total_inventory integer,
    total_reserved  integer
);

alter table room_type_inventory
    add constraint room_type_inventory_pk
        primary key (hotel_id, room_type_id, date);

INSERT INTO room_type_inventory (hotel_id, room_type_id, date, total_inventory, total_reserved)
VALUES
    (100, 1, '2023-10-12', 2, 0),
    (100, 1, '2023-10-13', 2, 0);
```

TODO: neden room değilde room_type_id tutuyoruz?
TODO: neden date'i teker teker tuttuk

Aynı istek 2.defa geldiğinde reservation id primary key olarak tuttuğumuz için
pq: duplicate key value violates unique constraint "reservation_pk" hatası ile double reservation'u önledik.

Pessimistic Lock için (FOR UPDATE için)
    - https://linuxhint.com/select-update-postgres/   
    - https://www.postgresql.org/docs/current/explicit-locking.html#LOCKING-ROWS

### Altering Room Type Inventory Table For Optimistic Locking
```sql
alter table room_type_inventory
add version integer;
```