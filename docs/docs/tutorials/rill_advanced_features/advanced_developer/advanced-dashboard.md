---
title: "Dashboard with more functionality"
description:  Further build on project
sidebar_label: "Advanced Dashboarding"
sidebar_position: 16
---

## Let's make a new dashboard

As we have learned in the previous course, we will need to setup the dashboard based on the new column names. 
Let's create a new dashboard via the UI. It should be named `dashboard_1.yaml`. Let's copy the contents from our old dashboard and make some changes.

First, we will want to change the `model` value to the new model name `advaned_commits___model`

Add two new dimensions: `directory path` and `commit_msg`.

Add four new measures: `SUM(total_line_changes)`, `SUM(net_line_changes)`, `SUM(num_commits)` and lastly let's create a percentage Code Deletion measure.

On the `SUM(net_line_changes)` measure, add the following `name: net_line_changes`. While name is not required, this can be used by other components for reference, which will be discussed later.

### Creating a measure in the metric-view
Like the SQL Model, our dashboards also use the same OLAP engine and you can use aggregates or expressions to create new metrics. In our case, since we have the added_lines and delete_lines measures, we can create a percentage of lines deleted measure.

```
 SUM(deleted_lines) / (SUM(added_lines) + SUM(deleted_lines))
```
:::tip
When to create measures in the SQL Model layer vs the metric-view layer?
It depends.

Depending on the size of data, type of measure, and what you are caluclating, you can choose either. Sometimes it would be better if you are dealing with a lot of data to front load the calculation on the SQL level so your dashboards load faster. However, the way OLAP engines work (linke avg of avg article), you might get incorrect data by doing certain calculations in the SQL level. You'll have to test and see which works for you!
:::

After making the above changes, you should be able to view your new dashboard!

![img](/img/tutorials/204/advanced-dashboard.png)

<details>
  <summary> Example Working Dashboard</summary>
  


```yaml
# Dashboard YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/dashboards

type: metrics_view

title: "My Advanced Tutorial Project"
#table: example_table # Choose a table to underpin your dashboard
model: advanced_commits___model

timeseries: author_date # Select an actual timestamp column (if any) from your table

dimensions: 
- column: directory_path
  label: "The directory"
  description: "The directory path"

- column: filename
  label: "The filename"
  description: "The name of the modified filename"

- column: author_name
  label: "The Author's Name"
  description: "The name of the author of the commit"

- column: commit_msg
  label: "The commit message"
  description: "The commit description attached."

measures:
- expression: "SUM(total_line_changes)"
  label: "Total number of Lines changed"
  description: "the total number of lines changes, addition and deletion"

- name: net_line_changes
  expression: "SUM(net_line_changes)"
  label: "Net number of Lines changed"
  description: "the total net number of lines changes"

- expression: "SUM(num_commits)"
  label: "Number of Commits"
  description: "The total number of commits"

- expression: "(SUM(deleted_lines)/(SUM(deleted_lines)+SUM(added_lines)))"
  label: "Code Deletion Percent %"
  description: "The percent of code deletion"
  format_preset: percentage
```

</details>


import DocsRating from '@site/src/components/DocsRating';

---
<DocsRating />