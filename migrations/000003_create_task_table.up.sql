CREATE TABLE task (
    TSK_ID serial not null,
    OPEN_D bigint,
    CLOSE_D bigint,
    CLOSED bool not null default false,
    META text not null,
    EMP_ID integer not null references employee(EMP_ID)
);