-- migrations/001_create_table.up.sql
CREATE TABLE IF NOT EXISTS Type_Display(
	ID_Type_Display serial not null constraint PK_Type_Display primary key,
	Name_Diagonal decimal(4, 1) not null,
	Name_Resolution varchar(70) not null,
	Type_Type varchar(70) not null,
	Type_Gsync boolean not null
);

CREATE TABLE IF NOT EXISTS Type_Monitor(
	ID_Type_Monitor serial not null constraint PK_Type_Monitor primary key,
	Name_Voltage decimal(4, 1) not null,
	Name_Gsync_Prem boolean not null,
	Name_Curved boolean not null,
	Type_Display_ID int not null references Type_Display(ID_Type_Display)
);

CREATE TABLE IF NOT EXISTS Type_Users(
	ID_Type_Users serial not null constraint PK_Type_Users primary key,
	Name_Username text not null constraint UQ_Type_Users unique,
	Name_Password text not null,
	Name_Email text not null,
	Name_Is_Admin boolean not null
);
