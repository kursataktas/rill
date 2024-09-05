---
title: Charts YAML
sidebar_label: Charts YAML
sidebar_position: 32
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

A Component consists of a data resolver (optional) and a output renderer, either Rill provided or a custom visualization defined as a Vega Lite specification.

## Properties

_**`type`**_ — Refers to the resource type and must be `component` _(required)_.

_**`data`**_ - A data resolver, either `metrics_sql`, `sql` or `api`. See the examples section for more detailed usage. _(required)_.

One of following output renderers:

_**`kpi`**_ - KPI object
    - `metric_view` - The metrics view to fetch data from _(required)_.
    - `time_range` - The time range, The value must be either a [valid ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations) or one of the [Rill ISO 8601 extensions](https://docs.rilldata.com/reference/rill-iso-extensions#extensions) _(required)_.
    - `measure` - A named measure from the `metrics_view` _(required)_.
    - `comparison_range`: the comparison time range, same requirements as `time_range`

_**`line_chart`**_ - Line Chart
    - `x` - X-axis values based off `metrics_sql` columns _(required)_
    - `y` - Y-axis values based off `metrics_sql` columns _(required)_
    - `xLabel` - A string for X-axis label
    - `yLabel` - A string for Y-axis label
    - `color` - A string to set the color, any valid html color string.

_**`bar_chart`**_ - Bar Chart
    - `x` - X-axis values based off `metrics_sql` columns _(required)_
    - `y` - Y-axis values based off `metrics_sql` columns _(required)_
    - `xLabel` - A string for X-axis label
    - `yLabel` - A string for Y-axis label
    - `color` - A string to set the color, any valid html color string.

_**`stacked_bar_chart`**_ - Stacked Bar Chart
    - `x` - X-axis values based off `metrics_sql` columns _(required)_
    - `y` - Y-axis values based off `metrics_sql` columns _(required)_
    - `xLabel` - A string for X-axis label
    - `yLabel` - A string for Y-axis label
    - `color` - A string to set the color, any valid html color string.

_*`markdown`*_ - Text Markdown
    - `content` - A markdown string that produces text.
    - `css` - A Object with key-value pairs of css properties.

_*`table`*_ - A Table object
    - `metric_view` - The metrics view to fetch data from _(required)_.
    - `col_dimensions` - List of named dimensions _(required)_
    - `row_dimensions` - List of named dimensions _(required)_
    - `measures` - List of named measure, will be used as columns _(required)_
    - `time_range` - The time range, The value must be either a [valid ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations) or one of the [Rill ISO 8601 extensions](https://docs.rilldata.com/reference/rill-iso-extensions#extensions) _(required)_.

_**`vega_lite`**_ - For any non-Rill based template charts, you need to define the vega_lite component.
    - `data` -  `{"name": "table"}` setting the data to the SQL query we defined under `data`  _(required)_
    - `mark` - defines the [type of chart](https://vega.github.io/vega-lite/docs/mark.html)  _(required)_
    - `encoding` - [Encoding the data](https://vega.github.io/vega-lite/docs/encoding.html) with visual properties   _(required)_
        - `x`, `y`, etc.
    - For more additional parameters to add, please refer to [Vega Lite's documentation](https://vega.github.io/vega-lite/docs/).


## Examples
### Rill Authored Charts
<Tabs>
<TabItem value="Bar" label="Bar Chart " default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts

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
</TabItem>

<TabItem value="Stacked Bar" label="Stacked Bar Chart " default>

```yaml

```

</TabItem>

<TabItem value="Line" label="Line Chart " default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts

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
</TabItem>

<TabItem value="Area" label="Area Chart " default>

```yaml

```
</TabItem>

<TabItem value="Scatter" label="Scatter Chart " default>

```yaml

```
</TabItem>

</Tabs>

### Rill Authored Others

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
![img](/img/build/canvasdashboard/KPI.png)
</TabItem>
<TabItem value="Table" label="Table" default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

table:
  measures:
    - net_line_changes
  metric_view: "dashboard_1"
  time_range: "P3M"
  #comparison_range: "P3M"

  row_dimensions:
    - author_name
  col_dimensions:
    - filename
```
</TabItem>

<TabItem value="Pivot Table" label="Pivot Table" default>

```yaml

```
</TabItem>

<TabItem value="Markdown" label="Markdown" default>

```yaml

```
</TabItem>

<TabItem value="Image" label="Image" default>

```yaml

```
</TabItem>

<TabItem value="Select" label="Select" default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

select:
  valueField: "1"
  label: "label"
  labelField: "labelField"
  placeholder: "Test"
```
</TabItem>

<TabItem value="Switch" label="Switch" default>

```yaml
# Chart YAML
# Reference documentation: https://docs.rilldata.com/reference/project-files/charts
    
type: component

switch:
  label: 'A switch'
  value: "1"
```
</TabItem>

</Tabs>

### Rill Authored Map Based

<Tabs>
<TabItem value="Layer Map" label="Layer Map" default>

```yaml

```
</TabItem>

<TabItem value="Choropleth" label="Choropleth Charts" default>

```yaml

```
</TabItem>


</Tabs>


### Vega-Lite Charts
<Tabs>

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

<img src = '/img/build/canvasdashboard/bar.png' class='rounded-gif' />
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

<img src = '/img/build/canvasdashboard/scatter.png' class='rounded-gif' />
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
<img src = '/img/build/canvasdashboard/line.png' class='rounded-gif' />
<br />

</TabItem>

</Tabs>

