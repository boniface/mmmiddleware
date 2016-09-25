create table cdn(
id text,
company text,
campaign text,
email text,
department text,
inref text,
extref text,
docType text,
docSize INT,
contentType text,
date text,
time text,
mongoId text,
filename text,
status text,
ref text,
reason text,
comment text,
PRIMARY KEY (company,mongoId)
);

create table cdnCampaign(
id text,
company text,
campaign text,
email text,
department text,
inref text,
extref text,
docType text,
docSize INT,
contentType text,
date text,
time text,
mongoId text,
filename text,
status text,
ref text,
reason text,
comment text,
PRIMARY KEY (company,campaign,id)
);
create table cdnEmail(
id text,
company text,
campaign text,
email text,
department text,
inref text,
extref text,
docType text,
docSize INT,
contentType text,
date text,
time text,
mongoId text,
filename text,
status text,
ref text,
reason text,
comment text,
PRIMARY KEY (company,email,id)
);
create table cdnDepartment(
id text,
company text,
campaign text,
email text,
department text,
inref text,
extref text,
docType text,
docSize INT,
contentType text,
date text,
time text,
mongoId text,
filename text,
status text,
ref text,
reason text,
comment text,
PRIMARY KEY (company,department,id)
);

create table cdnDocType(
id text,
company text,
campaign text,
email text,
department text,
inref text,
extref text,
docType text,
docSize INT,
contentType text,
date text,
time text,
mongoId text,
filename text,
status text,
ref text,
reason text,
comment text,
PRIMARY KEY (company,docType,id)
);

create table cdnDeleted(
id text,
company text,
email text,
mongoId text,
date text,
time text,
reason text,
comment text,
PRIMARY KEY (company,mongoId,id)
);
create table cdnDeletedDate(
id text,
company text,
email text,
mongoId text,
date text,
time text,
reason text,
comment text,
PRIMARY KEY (company,date,mongoId,id)
);

drop table cdn;
drop table cdncampaign;
drop table cdncompany;
drop table cdnDepartment;
drop table cdnEmail;
drop table cdnDocType;
