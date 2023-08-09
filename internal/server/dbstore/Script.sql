create table if not exists metrics(
	ID	varchar(50) not null,
	MTYPE	Varchar(10) not null,
	GAUGE	double precision,
	COUNTER INTEGER,
	PRIMARY KEY (ID, mtype)
);

create or replace function get_gauge( varchar) returns double precision as '
	select gauge  from metrics m
		where m.id = $1;'
	language sql;
end;

create or replace function get_counter( varchar) returns integer as '
	select counter  from metrics m
		where m.id = $1;'
	language sql;
end;

create or replace function set_gauge( p_id varchar, p_gauge double precision ) returns double precision as $$
	INSERT INTO metrics (id, mtype, gauge)
    VALUES (p_id, 'gauge', p_gauge)
    ON CONFLICT (id, mtype) DO UPDATE SET gauge = EXCLUDED.gauge 
   returning gauge;
	$$ language sql;
end;

create or replace function set_counter( p_id varchar, p_counter integer ) returns integer as $$
	INSERT INTO metrics (id, mtype, counter)
    VALUES (p_id, 'counter', p_counter)
    ON CONFLICT (id, mtype) DO UPDATE SET counter = EXCLUDED.counter 
   returning counter;
	$$ language sql;
end;
