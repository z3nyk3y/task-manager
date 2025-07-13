CREATE TABLE public.dict_task_status (
	id varchar(10) PRIMARY KEY
);

INSERT INTO public.dict_task_status (id)
VALUES ('NEW'), ('PROCESSING'), ('PROCESSED');


CREATE TABLE public.tasks (
	id serial8 PRIMARY KEY,
	status_id varchar(10) REFERENCES dict_task_status(id),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
	new.updated_at = now();

	RETURN new;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER set_updated_at
BEFORE UPDATE ON public.tasks
FOR EACH ROW
EXECUTE FUNCTION public.update_updated_at_column();
