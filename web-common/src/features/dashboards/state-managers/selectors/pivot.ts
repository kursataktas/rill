import type { DashboardDataSources } from "./types";
import { PivotChipType } from "../../pivot/types";

export const pivotSelectors = {
  showPivot: ({ dashboard }: DashboardDataSources) => dashboard.pivot.active,
  rows: ({ dashboard }: DashboardDataSources) => dashboard.pivot.rows,
  columns: ({ dashboard }: DashboardDataSources) => dashboard.pivot.columns,
  // Temporary limit to have only 1 time pill
  hasTimePill: ({ dashboard }: DashboardDataSources) => {
    const columns = dashboard.pivot.columns;
    const rows = dashboard.pivot.rows;

    return Boolean(
      columns.dimension.find((c) => c.type === PivotChipType.Time) ||
        rows.dimension.find((r) => r.type === PivotChipType.Time),
    );
  },
  measures: ({ metricsSpecQueryResult, dashboard }: DashboardDataSources) => {
    const measures =
      metricsSpecQueryResult.data?.measures?.filter(
        (d) => d.name && dashboard.visibleMeasureKeys.has(d.name),
      ) ?? [];
    const columns = dashboard.pivot.columns;

    return measures
      .filter((m) => !columns.measure.find((c) => c.id === m.name))
      .map((measure) => ({
        id: measure.name || "Unknown",
        title: measure.label || measure.name || "Unknown",
        type: PivotChipType.Measure,
      }));
  },
  dimensions: ({ metricsSpecQueryResult, dashboard }: DashboardDataSources) => {
    {
      const dimensions =
        metricsSpecQueryResult.data?.dimensions?.filter(
          (d) => d.name && dashboard.visibleDimensionKeys.has(d.name),
        ) ?? [];
      const columns = dashboard.pivot.columns;
      const rows = dashboard.pivot.rows;

      return dimensions
        .filter((d) => {
          return !(
            columns.dimension.find((c) => c.id === d.name) ||
            rows.dimension.find((r) => r.id === d.name)
          );
        })
        .map((dimension) => ({
          id: dimension.name || dimension.column || "Unknown",
          title:
            dimension.label || dimension.name || dimension.column || "Unknown",
          type: PivotChipType.Dimension,
        }));
    }
  },
};
