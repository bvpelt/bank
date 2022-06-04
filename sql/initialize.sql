create table target (
    id bigserial,
    name text not null,
    description text,
    primary key (id),
    unique (name)
);

create table account (
    id bigserial,
    number text not null,
    description text,
    primary key (id),
    unique (number)
);

create table transaction (
    id bigserial,
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