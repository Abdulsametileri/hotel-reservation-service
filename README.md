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

### Altering Room Type Inventory Table For Optimistic Locking

```sql
alter table room_type_inventory
add version integer default 0;
```

### Adding reservation db
```sql
create database reservation;
```

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

### Adding inventory db
```sql
create database inventory;
```

```sql
create table room_type_inventory
(
    hotel_id        integer,
    room_type_id    integer,
    date            date,
    total_inventory integer,
    total_reserved  integer
);
```