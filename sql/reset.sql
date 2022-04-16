
delete from account;
delete from target;
delete from transaction;

ALTER SEQUENCE account_id_seq RESTART;
ALTER SEQUENCE target_id_seq RESTART;
ALTER SEQUENCE transaction_id_seq RESTART;
