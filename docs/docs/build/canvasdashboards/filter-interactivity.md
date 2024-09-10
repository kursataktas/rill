---
title: Variables
description: Utilize variables to communicate between components.
sidebar_label: Variables
sidebar_position: 10
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

## Variables - Communicate between components

Most dashboards requires some level of interactivity, such as filtering data. In Rill Canvas Dashboards this is achieved by utilizing variables.
Variables allows components to communicate with each other and pass values around within a Canvas Dashboard.

### Components - Input and Outputs
A Component can either be a producer of a value (refered to as output) or a consumer of a value (refered to as input).
A Canvas Dashboard can contain variables that a component can either write to or consume. A example scenario could be that you have a Canvas dashboard that has a variable named `country`. In your dashboard you have a Select box that produces the currently selected value in the listbox to the Canvas Dashboard. That variable is later passed down to components that accepts input on the same name and can use that value in a `where` clause to filter data for example.

Piecing all of this together allows the developer to set up filters that can span multiple different models and metrics views when there is a need for a shared dimension.



