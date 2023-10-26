drop trigger set_logs_updated_timestamp on logs;
drop trigger data_cleanup_trigger on logs;
drop function delete_old_data();

drop table if exists logs;