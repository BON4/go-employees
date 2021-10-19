CREATE TABLE employee (
      EMP_ID serial not null primary key,
      FNAME varchar (50) not null,
      LNAME varchar (50) not null,
      SAL numeric default 0
);