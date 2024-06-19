# Druid Task Metatdata Cleanup

This tool is meant to clean up task metadata for Apache Druid when using Postgres as the metadata database and S3 as the storage backend. To use the tool, you first must set up the configuration, which looks like this YAML file.

```yaml
database:
  user: druid      # Postgres admin username
  password: diurd  # Postgres admin password
  host: localhost  # Postgres host address
  port: 5432       # Postgres host port
  db: druid        # Postgres metatdata database for Druid cluster

s3:
  bucket: my-druid-bucket  # S3 bucket name for Druid metadata
  keyBase: indexing-logs   # Path in bucket to indexing logs 
  batchSize: 1000          # How many objects to delete from S3 at a time

endTime: "2021-09-01"  # The end created_date to delete metadata up to
```

Replace the fields with the values for your deployment of Druid.

Build the tool with `go build`, then you can run it (or just `go run .`). Examples for using the tool:

```bash
./druid-metadata-cleanup  # Will only show the number of records in Postgres to delete and the number of S3 objects to delete. Will not actually delete them
./druid-metadata-cleanup --show-tasks # Run in preview mode, but also print all the task IDs from Postgres
./druid-metadata-cleanup --config-file my-config.yaml # Run the tool, but with a different config file than default (default is config.yaml)
./druid-metadata-cleanup --show-tasks --config-file my-config.yaml --delete  # Run the tool and enable the deletion of Postgres records and S3 objects
```
