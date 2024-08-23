## notemplate

a program for creating notes

```sh
go run entry.go
# create default entry in notemplate/
# => notemplates/2024-08-23-entry-0.toml

go run entry.go _ withsuffix
# create default entry with suffix in notemplate/
# => notemplates/2024-08-23-entry-1-withsuffix.toml

go run entry.go job_applications
# create job_application template entry in job_applications/
# => job_applications/2024-08-23-entry-0.toml

go run entry.go job_applications withsuffix
# create job_application template entry with suffix in job_applications/
# => job_applications/2024-08-23-entry-withsuffix-1.toml
```