---
title: "Splits with Cloud Storage"
description:  "Getting Started with Splits"
sidebar_label: "Cloud Storage: Splits and Incremental Models"
sidebar_position: 12
---

Another advanced concept within Rill is using [Incremental Models](https://docs.rilldata.com/build/advancedmodels/incremental). To understand incremental models, we will also need to discuss [splits](https://docs.rilldata.com/build/advancedmodels/splits). 


## Understanding Splits in Models

Here’s how it works at a high level:

- **Split Definition**: Each row from the result set becomes one "split". The model processes each split separately.
- **Execution Strategy**:
  - **First Split**: Runs without incremental processing.
  - **Subsequent Splits**: Run incrementally, following the output connector's `incremental_strategy` (either append or merge for SQL connectors).

### Let's create a basic split model.
In the previous courses, we used a GCS connection to import ClickHouse's repository commit history. In this guide, we will use S3. The format of the files are the same, you just need to change `gs` to `s3`


```
s3://rilldata-public/github-analytics/Clickhouse/*/*/commits_*.parquet
```
1. Create a YAML file: `S3-splits-tutorial.yaml`

2. Use `glob:` resolver to load files from S3
```yaml
splits:
  glob:
    connector: s3
    path: s3://rilldata-public/github-analytics/Clickhouse/*/*/commits_*.parquet

```
3. Set the SQL statement to user the URI.
```yaml
sql: SELECT * FROM read_parquet('{{ .split.uri }}')
```

### Handling errors in splits
If you see any errors in the UI regarding split, you may need to check the status. You can do this via the CLI running:
```bash
rill project splits --<model_name> --local
```

Once completed you should see the following:

![img](/img/tutorials/302/splits.png)

### Refreshing Splits 

Let's say a specific split in your model had some formatting issues. After fixing the data, you would need to find the key for the split and run `rill project splits --<model_name> --local`.  Once found, you can run the following command that will only refresh the specific split, instead of the whole model.

```bash
rill project refresh --model <model_name> --split <split_key>
```


## What is Incremental Modeling?
Once splits are setup, you can use incremental modeling to load only new data when refreshing a dataset. This becomes important when your data is large and it does not make sense to reload all the data when trying to ingest new data.

### Let's modify the split model to add incremental modeling.

1. Set `incremental` to true

2. You can manually setup a `splits_watermark` but since our data is using the `glob` key, it is automatically set to the `updated_on` field. 

3. Let's set up a `refresh` based on `cron` that runs daily at 8AM UTC.
```
refresh:
    cron: "0 8 * * *"
```

Once Rill ingests the data, your UI should looks something like this: 
![img](/img/tutorials/302/incremental.png)

Your YAML should look like the following:

```yaml
type: model 

incremental: true 
refresh:
  cron: 0 8 * * *

splits:
  glob:
    connector: s3
    path: s3://rilldata-public/github-analytics/Clickhouse/*/*/commits_*.parquet

sql: SELECT * FROM read_parquet('{{ .split.uri }}')

```

You now have a working incremental model that refreshed new data based on the `updated_on` key at 8AM UTC everyday. Along with writing to the default OLAP engine, DuckDB, we have also added some features to use staging tables for connectors that do not have direct read/write capabilities.

import DocsRating from '@site/src/components/DocsRating';

---
<DocsRating />