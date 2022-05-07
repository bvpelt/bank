
delete from transaction;
delete from account;
delete from target;

ALTER SEQUENCE account_id_seq RESTART;
ALTER SEQUENCE target_id_seq RESTART;
ALTER SEQUENCE transaction_id_seq RESTART;
