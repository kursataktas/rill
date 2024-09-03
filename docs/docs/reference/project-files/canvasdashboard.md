---
title: Canvas Dashboard YAML
sidebar_label: Canvas Dashboard YAML
sidebar_position: 31
---

## Properties

**`type`** - Refers to the resource type and must be `dashboard` _(required)_.

**`title`** - (String) The title of the Canvas Dashboard

**`columns`** - (Int) Number of columns in the grid

**`gap`** - (Int) Gap between cells in the grid

**`items`** List of components in the Canvas dashboard
    - **`component`** — Either a reference to a named component file or a inline decleration of a component. _(required)_
    - **`x`** — x starting position in the grid. _(required)_
    - **`y`** — y starting position in the grid. _(required)_
    - **`width`** — number of columns to span. _(required)_
    - **`height`** — number of rows to span. _(required)_ 


## Example

```yaml
type: dashboard

title: My Dashboard
columns: 32
gap: 2

items:
  - component: mytablecomponent
    x: 2
    y: 2
    width: 10
    height: 3
  - component:
      markdown:
        content: "This is my text"
    x: 2
    y: 6
    width: 3
    height: 2
```