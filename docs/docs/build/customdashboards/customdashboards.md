---
title: Create Custom Dashboards
description: Build custom dashboards, combining multiple metric views and free form visualizations
sidebar_label: Create Custom Dashboards
sidebar_position: 00
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';


In Rill Custom Dashboards allows you to build more traditional dashboards that combines data from multiple different metric views and gives you a higher degree of freedom in terms of design and layout.

## Getting started
> not sure if/when we will remove the feature flag, but will add this for now.

In order to enable custom dashboarding in your environment, you will need to enable the feature flag.

```
features:
 - customDashboards
```

Once enabled, you will see a few more items populate in the `Add` dropdown.

![img](/img/build/customdashboard/add-custom-dashboard.png)


## Anatomy of a Custom Dashboard

When getting started on a new custom dashboard, you will be introduced to this following page:

![img](/img/build/customdashboard/custom-dashboard.png)

There are three views for a custom dashboard. 
- Code View
- Split View
- Viz View

You can either add charts that you've already made or create components directly within the dashboard. If you want to add a chart you've already created you can select the `+ Add chart` in the UI or start writing YAML to bring in the component into the custom dashboard.

:::note
The syntax to create your component may vary from an independent chart vs. creating it in a dashboard. We'll go into more details on this topic later on.

<Tabs>
<TabItem value="file" label="Rill KPI Chart as a file, imported into dashboard" default>
Chart.yaml:
```yaml
type: component

kpi:
  metric_view: dashboard_1
  time_range: P1W
  measure: measure_2
  comparison_range: P1W
```

Dashboard.yaml:
```yaml
  - component: net_line_kpi
    height: 2
    width: 4
    x: 0
    y: 1
```
</TabItem>
<TabItem value="dashboard" label="KPI Chart created directly on a dashboard">

```yaml
  - component:
      kpi:
        metric_view: dashboard_1
        time_range: P1W
        measure: measure_2
        comparison_range: P1W
    width: 3
    height: 1
    x: 0
    y: 0
```
</TabItem>

</Tabs>
:::


## Understanding Components

There are a few different types of charts that can be added to a custom dashboard and this determines the required components.

### Rill KPI Templates
If you are using a Rill template, you can call a **metric-view** directly and use the already defined components in the metric-view. In this case, you will not need to define a Vega Lite component and view specification to build the chart.

### Vega Lite charts
Using [Vega Lite's chart creating capablities](https://vega.github.io/vega-lite/docs/spec.html), we allow you to customize a chart to whatever your needs are. Within the chart YAML file, you will need to define the `type`, `data` and `vega_lite` component.

We will go into more details in the next page, [Create Components](components/).

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
<TabItem value="Bar" label="Vega Lite -  Bar Charts">

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

<TabItem value="Scatter" label="Vega Lite -  Scatter Charts">

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

<TabItem value="Line" label="Vega Lite -  Line Charts">

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