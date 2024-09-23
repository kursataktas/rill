---
title: Components
description: Components make up your Canvas dashboards
sidebar_label: Create Components
sidebar_position: 00
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


In Rill Canvas Dashboasrds allows you to build more traditional dashboards that combines data from multiple different metric views and gives you a higher degree of freedom in terms of design and layout.

## Anatomy of a Component
Every component consists of four main sections:

**`data`** - This section defines which data should be exposed to the component. Various different data resolvers are available such as `metrics_sql`, `api` and `sql`

**`input`** - Defines the various input parameters the component expects. Either supplied via a Canvas Dashboard or from other external sources such as in a embedded scenario.

**`output`** - Defines the various outputs that the component can produce. Typically they would map to a Canvas Dashboard variable.

**`renderer (kpi/table/bar/line/vega_lite)`** - The visual output that gets rendered onto the screen. Either a template that Rill has created or your own custom visualization provided via vega.


## Rill Authored Components
Available component types:


<div style={{ display: 'flex', justifyContent: 'space-between' }}>

  <div style={{ flex: '1', padding: '10px' }}>
    <!-- Column 1 content goes here -->
    ### Charts
      - Bar Chart
      - Stacked Bar Chart
      - Line Chart
      - KPI
      - Table
  </div>

  <div style={{ flex: '1', padding: '10px' }}>
    <!-- Column 2 content goes here -->
    ### Others
      - Markdown
      - Image
      - Select
      - Switch
  </div>

</div>


You will need to define the `data` component using a `sql` statement (from table) or `metric_sql` statement (from dashboard).
```yaml
data:
    [metric_sql]: |
```

Once this is done, you can set the graph type with [`bar_chart` or `line_chart`] and set the `x` and `y` axis.
```yaml
bar_chart:
  x: x-axis column
  y: y-axis column
```

## Building your own components

### Vega Lite 

Vega-Lite is a high-level grammar of interactive graphics, built on top of the Vega visualization framework. It allows users to create a wide range of visualizations such as bar charts, line graphs, scatter plots, and more with concise JSON syntax. Vega-Lite focuses on simplicity and expressiveness, enabling the creation of complex visualizations with minimal code. It supports data transformations, filtering, and aggregation, making it ideal for exploring, analyzing, and presenting data interactively.

Rill allows users to write their own chart specificiations in Vega-Lite that nativly integrates with Rills data engine. Before we can build our own component, we will need some data to work with. The first step for every component would be to declare a data section that defines which data that should be available to the component.

```
data:
    [metric_sql]: |
```
You can either use vanilla SQL based on your underlying models or sources via `sql:` to select the data that you need. Or, you can use `metric_sql` which allows you to selec data directly from your dashboard's metrics. 

:::note
Keep in mind that vega-lite is not well-suited for rendering large number of elements. If you are having issues rendering the chart, we recommened using appending a `limit XXXX` to your SQL.
:::

Now that we have data, we can build out our Vega Lite component.



## Examples

Check out our [references](../../reference/project-files/components.md#examples) for more examples!
<Tabs>
<TabItem value="KPI" label="KPI Chart " default>

```yaml
type: component

kpi:
  metric_view: dashboard_1
  time_range: P1W
  measure: measure_2
  comparison_range: P1W

```

<img src = '/img/build/canvasdashboard/kpi.png' class='rounded-gif' />
<br />
</TabItem>

<TabItem value="Rill_Chart" label="Rill Authored Chart " default>

```yaml
type: component

data:
  metrics_sql: |
    select 
      measure_0,
      date_trunc('day', author_date) as date 
    from dashboard_1
    where author_date > '2024-07-14 00:00:00 Z'

line_chart:
  x: date
  y: measure_0
```

<img src = '/img/build/canvasdashboard/rill-chart.png' class='rounded-gif' />
<br />
</TabItem>

<TabItem value="Bar" label="Vega -  Bar Charts">

```yaml 
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

<img src = '/img/build/canvasdashboard/bar.png' class='rounded-gif' />
<br />
</TabItem>

<TabItem value="Line" label="Vega - Line Charts">

```yaml
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
<img src = '/img/build/canvasdashboard/line.png' class='rounded-gif' />
<br />

</TabItem>

</Tabs>