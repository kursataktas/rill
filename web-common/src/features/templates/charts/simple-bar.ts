import { ChartField } from "@rilldata/web-common/features/templates/charts/build-template";
import { singleLayerBaseSpec } from "./utils";
import { ScrubBoxColor } from "@rilldata/web-common/features/dashboards/time-series/chart-colors";

export function buildSimpleBar(
  timeField: ChartField,
  quantitativeField: ChartField,
) {
  const baseSpec = singleLayerBaseSpec();

  baseSpec.mark = {
    type: "bar",
    width: { band: 0.75 },
    clip: true,
  };

  baseSpec.encoding = {
    x: {
      field: timeField.name,
      type: "temporal",
      bandPosition: 0,
    },
    y: { field: quantitativeField.name, type: "quantitative" },
    opacity: {
      condition: [
        {
          param: "hover",
          empty: false,
          value: 1,
        },
        {
          param: "brush",
          empty: false,
          value: 1,
        },
      ],
      value: 0.8,
    },
    // TODO: configure or disable tooltip while scrubbing for the time being
    // https://vega.github.io/vega-lite/docs/tooltip.html#disable-tooltips
    // TODO: can add a `disableTooltip` flag to buildSimpleBar
    // TODO: do we want to turn this into multiLayerSpec to disable tooltip on interval?
    tooltip: [
      {
        field: timeField.tooltipName ? timeField.tooltipName : timeField.name,
        type: "temporal",
        title: "Time",
        format: "%b %d, %Y %H:%M",
      },
      {
        title: quantitativeField.label,
        field: quantitativeField.name,
        type: "quantitative",
        formatType: quantitativeField.formatterFunction || "number",
      },
    ],
  };

  baseSpec.params = [
    {
      name: "hover",
      select: {
        type: "point",
        on: "pointerover",
        clear: "pointerout",
        encodings: ["x"],
      },
    },
    {
      name: "brush",
      select: {
        type: "interval",
        encodings: ["x"],
        mark: {
          fill: ScrubBoxColor,
          fillOpacity: 0.2,
          stroke: ScrubBoxColor,
          strokeWidth: 1,
          strokeOpacity: 0.8,
        },
        // TODO: create event stream to clear brush on escape key
        // https://vega.github.io/vega-lite-v4/docs/clear.html
        clear: "dblclick",
      },
    },
  ];

  return baseSpec;
}
