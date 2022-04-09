create table target (
    id bigint not null,
    name text not null,
    description text,
    primary key (id)
);

create table account (
    id bigint not null,
    number text not null,
    description text,
    primary key (id)
);

create table transaction (
    id bigint not null,
    from_account bigint not null,
    to_account bigint not null,
    target bigint not null,
    amount bigint, -- referenced to from_account
    description text,
    primary key (id),
    foreign key (from_account) references account (id),
    foreign key (to_account) references account (id),
    foreign key (target) references target (id)
);