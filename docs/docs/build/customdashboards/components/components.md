---
title: Components
description: Components make up your custom dashboards
sidebar_label: Create Components
sidebar_position: 00
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


In Rill Custom Dashboards allows you to build more traditional dashboards that combines data from multiple different metric views and gives you a higher degree of freedom in terms of design and layout.

## Rill Authored Visualization Components

## Rill Authored Interaction Components

## Building your own component using Vega Lite
Before we can build our own component with Vega Lite, we will need some data to work with.

```
data:
    [sql/metric_sql]: |
```
You can either use vanilla SQL based on your underlying models or sources via `sql:` to select the data that you need. Or, you can use `metric_sql` which allows you to selec data directly from your dashboard's metrics. 

:::note
Keep in mind that vega-lite is not well-suited for rendering large number of elements. If you are having issues rendering the chart, we recommened using appending a `limit XXXX` to your SQL.
:::

Now that we have data, we can build out our Vega Lite component.

### Vega Lite stuff

## Examples

<Tabs>
<TabItem value="KPI" label="KPI Chart " default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

kpi:
  metric_view: dashboard_1 #name of dashboard
  time_range: P1W
  measure: measure_2  #retrieved from the dashboard top down [0:]
  comparison_range: P1W

```

![img](/img/build/customdashboard/kpi.png)
</TabItem>
<TabItem value="Bar" label="Vega_lite -  Bar Charts">

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

data:
  sql: |
      SELECT
          category,
          value
      FROM
          (VALUES
              ('Adidas', 50),
              ('Nike', 80),
              ('FILA', 45),
              ('Converse', 45),
          ) AS data(category, value)


vega_lite: |
  {
    "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
    "data": { "name": "table" },
    "mark": "bar",
    "width": "container",
    "height": 500,
    "encoding": {
      "x": {
        "field": "category",
        "type": "ordinal",
        "axis": { "title": "Category" }
      },
      "y": {
        "field": "value",
        "type": "quantitative",
        "axis": { "title": "Value" }
      }
    }
  }
```

<img src = '/img/build/customdashboard/bar.png' class='rounded-gif' />
<br />
</TabItem>

<TabItem value="Scatter" label="Vega_lite -  Scatter Charts">

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

data:
  sql: |
    SELECT * FROM (VALUES 
        ('Point A', 1, 2),
        ('Point B', 2, 3),
        ('Point C', 3, 6),
        ('Point D', 4, 8),
        ('Point E', 5, 5),
        ('Point F', 6, 7),
        ('Point G', 7, 10),
        ('Point H', 8, 6),
        ('Point I', 9, 9),
        ('Point J', 10, 12)
    ) AS t(name, x_value, y_value)

vega_lite: |
  {
    "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
    "data": { "name": "table" },
    "mark": "point",
    "width": "container",
    "height": 500,
    "encoding": {
      "x": {
        "field": "x_value",
        "type": "quantitative",
        "axis": { "title": "Point"}
      },
      "y": {
        "field": "y_value",
        "type": "quantitative",
        "axis": { "title": "Point", "orient": "left"  }
      },
        "color": {
          "field": "name",
          "type": "nominal",
          "legend": {"title": "Point Name"}
        },
        "tooltip": [
          {"field": "name", "type": "nominal", "title": "Point Name"},
          {"field": "x_value", "type": "quantitative", "title": "X Value"},
          {"field": "y_value", "type": "quantitative", "title": "Y Value"}
        ]
    }
  }
```

<img src = '/img/build/customdashboard/scatter.png' class='rounded-gif' />
<br />
</TabItem>

<TabItem value="Line" label="Vega_lite -  Line Charts">

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

data:
  sql: |
    SELECT * FROM (VALUES 
      ('Monday', 300),
      ('Tuesday', 150),
      ('Wednesday', 200),
      ('Thursday', 400),
      ('Friday', 650),
      ('Saturday', 575),
      ('Sunday', 500)
    ) AS t(day_of_week, revenue)

vega_lite: |
  {
    "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
    "data": { "name": "table" },
    "mark": "line",
    "width": "container",
    "height": 500,
    "encoding": {
      "x": {
        "field": "day_of_week",
        "type": "ordinal",
        "axis": { "title": "Day of the Week" },
        "sort": [
          "Monday",
          "Tuesday",
          "Wednesday",
          "Thursday",
          "Friday",
          "Saturday",
          "Sunday"
        ]
      },
      "y": {
        "field": "revenue",
        "type": "quantitative",
        "axis": { "title": "Revenue" }
      }
    }
  }
```
<img src = '/img/build/customdashboard/line.png' class='rounded-gif' />
<br />

</TabItem>

</Tabs>