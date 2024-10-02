import { splitWhereFilter } from "@rilldata/web-common/features/dashboards/filters/measure-filters/measure-filter-utils";
import {
  PivotChipData,
  PivotChipType,
  PivotState,
} from "@rilldata/web-common/features/dashboards/pivot/types";
import { SortDirection } from "@rilldata/web-common/features/dashboards/proto-state/derived-types";
import { createAndExpression } from "@rilldata/web-common/features/dashboards/stores/filter-utils";
import type {
  DimensionThresholdFilter,
  MetricsExplorerEntity,
} from "@rilldata/web-common/features/dashboards/stores/metrics-explorer-entity";
import {
  TDDChart,
  TDDState,
} from "@rilldata/web-common/features/dashboards/time-dimension-details/types";
import {
  URLStateDefaultSortDirection,
  URLStateDefaultTimezone,
} from "@rilldata/web-common/features/dashboards/url-state/defaults";
import { convertFilterParamToExpression } from "@rilldata/web-common/features/dashboards/url-state/filters/converters";
import {
  FromURLParamTDDChartMap,
  FromURLParamTimeDimensionMap,
  FromURLParamViewMap,
  ToActivePageViewMap,
} from "@rilldata/web-common/features/dashboards/url-state/mappers";
import { getMapFromArray } from "@rilldata/web-common/lib/arrayUtils";
import { TIME_GRAIN } from "@rilldata/web-common/lib/time/config";
import { DashboardTimeControls } from "@rilldata/web-common/lib/time/types";
import {
  MetricsViewSpecDimensionV2,
  MetricsViewSpecMeasureV2,
  type V1ExplorePreset,
  V1ExploreSpec,
  V1ExploreWebView,
  V1Expression,
  V1MetricsViewSpec,
  V1Operation,
} from "@rilldata/web-common/runtime-client";

export function getMetricsExplorerFromUrl(
  searchParams: URLSearchParams,
  metricsView: V1MetricsViewSpec,
  explore: V1ExploreSpec,
  preset: V1ExplorePreset,
): { entity: Partial<MetricsExplorerEntity>; errors: Error[] } {
  // TODO: replace this with V1ExplorePreset once it is available on main
  const entity: Partial<MetricsExplorerEntity> = {};
  const errors: Error[] = [];

  const measures = getMapFromArray(
    metricsView.measures?.filter((m) => explore.measures?.includes(m.name!)) ??
      [],
    (m) => m.name!,
  );
  const dimensions = getMapFromArray(
    metricsView.dimensions?.filter((d) =>
      explore.dimensions?.includes(d.name!),
    ) ?? [],
    (d) => d.name!,
  );

  if (searchParams.has("vw")) {
    entity.activePage = Number(
      ToActivePageViewMap[
        FromURLParamViewMap[searchParams.get("vw") as string]
      ] ?? "0",
    );
  }

  if (searchParams.has("f")) {
    const {
      dimensionFilters,
      dimensionThresholdFilters,
      errors: filterErrors,
    } = fromFilterUrlParam(searchParams.get("f") as string);
    if (filterErrors) errors.push(...filterErrors);
    if (dimensionFilters) entity.whereFilter = dimensionFilters;
    if (dimensionThresholdFilters)
      entity.dimensionThresholdFilters = dimensionThresholdFilters;
  }

  const { entity: trEntity, errors: trErrors } = fromTimeRangesParams(
    searchParams,
    dimensions,
    preset,
  );
  Object.assign(entity, trEntity);
  errors.push(...trErrors);

  Object.assign(
    entity,
    fromOverviewUrlParams(searchParams, measures, dimensions, explore, preset),
  );

  entity.tdd = fromTimeDimensionUrlParams(searchParams, measures, preset);

  entity.pivot = fromPivotUrlParams(searchParams, measures, dimensions, preset);

  return { entity, errors };
}

function fromTimeRangesParams(
  searchParams: URLSearchParams,
  dimensions: Map<string, MetricsViewSpecDimensionV2>,
  preset: V1ExplorePreset,
) {
  const entity: Partial<MetricsExplorerEntity> = {};
  const errors: Error[] = [];

  const timeRange = preset.timeRange || searchParams.get("tr");
  if (timeRange) {
    const { timeRange: selectedTimeRange, error } =
      fromTimeRangeUrlParam(timeRange);
    if (error) errors.push(error);
    entity.selectedTimeRange = selectedTimeRange;
  }
  const timeZone = preset.timezone || searchParams.get("tz");
  if (timeZone) {
    entity.selectedTimezone = timeZone;
  } else {
    entity.selectedTimezone = URLStateDefaultTimezone;
  }

  const comparisonTimeRange =
    preset.compareTimeRange || searchParams.get("ctr");
  if (comparisonTimeRange) {
    const { timeRange, error } = fromTimeRangeUrlParam(comparisonTimeRange);
    if (error) errors.push(error);
    entity.selectedComparisonTimeRange = timeRange;
  }
  const comparisonDimension =
    preset.comparisonDimension || searchParams.get("cd");
  if (comparisonDimension && dimensions.has(comparisonDimension)) {
    entity.selectedComparisonDimension = comparisonDimension;
  }

  return { entity, errors };
}

function fromFilterUrlParam(filter: string): {
  dimensionFilters?: V1Expression;
  dimensionThresholdFilters?: DimensionThresholdFilter[];
  errors?: Error[];
} {
  try {
    let expr = convertFilterParamToExpression(filter);
    // if root is not AND/OR then add AND
    if (
      expr?.cond?.op !== V1Operation.OPERATION_AND &&
      expr?.cond?.op !== V1Operation.OPERATION_OR
    ) {
      expr = createAndExpression([expr]);
    }
    return splitWhereFilter(expr);
  } catch (e) {
    return { errors: [e] };
  }
}

function fromTimeRangeUrlParam(tr: string): {
  timeRange?: DashboardTimeControls;
  error?: Error;
} {
  // TODO: validation
  return {
    timeRange: {
      name: tr,
    } as DashboardTimeControls,
  };

  // return {
  //   error: new Error(`unknown time range: ${tr}`),
  // };
}

function fromOverviewUrlParams(
  searchParams: URLSearchParams,
  measures: Map<string, MetricsViewSpecMeasureV2>,
  dimensions: Map<string, MetricsViewSpecDimensionV2>,
  explore: V1ExploreSpec,
  preset: V1ExplorePreset,
) {
  const entity: Partial<MetricsExplorerEntity> = {};

  let selectedMeasures = preset.measures ?? explore.measures ?? [];
  if (searchParams.has("o.m")) {
    const mes = searchParams.get("o.m") as string;
    if (mes !== "*") {
      selectedMeasures = mes.split(",").filter((m) => measures.has(m));
    }
  }
  entity.allMeasuresVisible =
    selectedMeasures.length === explore.measures?.length;
  entity.visibleMeasureKeys = new Set(selectedMeasures);

  let selectedDimensions = preset.dimensions ?? explore.dimensions ?? [];
  if (searchParams.has("o.d")) {
    const dims = searchParams.get("o.d") as string;
    if (dims !== "*") {
      selectedDimensions = dims.split(",").filter((d) => dimensions.has(d));
    }
  }
  entity.allDimensionsVisible =
    selectedDimensions.length === explore.dimensions?.length;
  entity.visibleDimensionKeys = new Set(selectedDimensions);

  entity.leaderboardMeasureName =
    preset.overviewSortBy ?? preset.measures?.[0] ?? explore.measures?.[0];
  if (searchParams.has("o.sb")) {
    const sortBy = searchParams.get("o.sb") as string;
    if (measures.has(sortBy)) {
      entity.leaderboardMeasureName = sortBy;
    }
  }

  if (preset.overviewSortAsc !== undefined) {
    entity.sortDirection = preset.overviewSortAsc
      ? SortDirection.ASCENDING
      : SortDirection.DESCENDING;
  } else {
    entity.sortDirection = URLStateDefaultSortDirection;
  }
  if (searchParams.has("o.sd")) {
    const sortDir = searchParams.get("o.sd") as string;
    entity.sortDirection =
      sortDir === "ASC" ? SortDirection.ASCENDING : SortDirection.DESCENDING;
  }

  entity.selectedDimensionName = preset.overviewExpandedDimension ?? "";
  if (searchParams.has("o.ed")) {
    const dim = searchParams.get("o.ed") as string;
    if (dimensions.has(dim)) {
      entity.selectedDimensionName = dim;
    }
  }

  return entity;
}

function fromTimeDimensionUrlParams(
  searchParams: URLSearchParams,
  measures: Map<string, MetricsViewSpecMeasureV2>,
  preset: V1ExplorePreset,
) {
  let ttdMeasure = preset.timeDimensionMeasure;
  let ttdChartType = preset.timeDimensionChartType as TDDChart | undefined;
  let ttdPin: number | undefined; // TODO

  if (searchParams.has("tdd.m")) {
    const mes = searchParams.get("tdd.m") as string;
    if (measures.has(mes)) {
      ttdMeasure = mes;
    }
  }
  if (searchParams.has("tdd.ct")) {
    const ct = searchParams.get("tdd.ct") as string;
    if (ct in FromURLParamTDDChartMap) {
      ttdChartType = FromURLParamTDDChartMap[ct];
    }
  }
  if (searchParams.has("tdd.p")) {
    const pin = Number(searchParams.get("tdd.p") as string);
    if (!Number.isNaN(pin)) {
      ttdPin = pin;
    }
  }

  return <TDDState>{
    expandedMeasureName: ttdMeasure ?? "",
    chartType: ttdChartType ?? TDDChart.DEFAULT,
    pinIndex: ttdPin ?? -1,
  };
}

function fromPivotUrlParams(
  searchParams: URLSearchParams,
  measures: Map<string, MetricsViewSpecMeasureV2>,
  dimensions: Map<string, MetricsViewSpecDimensionV2>,
  preset: V1ExplorePreset,
): PivotState {
  const mapPivotEntry = (entry: string): PivotChipData | undefined => {
    if (entry in FromURLParamTimeDimensionMap) {
      const grain = FromURLParamTimeDimensionMap[entry];
      return {
        id: grain,
        title: TIME_GRAIN[grain]?.label,
        type: PivotChipType.Time,
      };
    }

    if (measures.has(entry)) {
      const m = measures.get(entry)!;
      return {
        id: entry,
        title: m.label || m.name || "Unknown",
        type: PivotChipType.Measure,
      };
    }

    if (dimensions.has(entry)) {
      const d = dimensions.get(entry)!;
      return {
        id: entry,
        title: d.label || d.name || "Unknown",
        type: PivotChipType.Dimension,
      };
    }

    return undefined;
  };

  const rowDimensions: PivotChipData[] = [];
  let pivotRows = preset.pivotRows;
  if (searchParams.has("p.r")) {
    pivotRows = (searchParams.get("p.r") as string).split(",");
  }
  if (pivotRows) {
    pivotRows.forEach((pivotRow) => {
      const chip = mapPivotEntry(pivotRow);
      if (!chip) return;
      rowDimensions.push(chip);
    });
  }

  const colMeasures: PivotChipData[] = [];
  const colDimensions: PivotChipData[] = [];
  let pivotCols = preset.pivotCols;
  if (searchParams.has("p.c")) {
    pivotCols = (searchParams.get("p.c") as string).split(",");
  }
  if (pivotCols) {
    pivotCols.forEach((pivotRow) => {
      const chip = mapPivotEntry(pivotRow);
      if (!chip) return;
      if (chip.type === PivotChipType.Measure) {
        colMeasures.push(chip);
      } else {
        colDimensions.push(chip);
      }
    });
  }

  return {
    active:
      searchParams.get("vw") === "pivot" ||
      preset.view === V1ExploreWebView.EXPLORE_ACTIVE_PAGE_PIVOT,
    rows: {
      dimension: rowDimensions,
    },
    columns: {
      measure: colMeasures,
      dimension: colDimensions,
    },
    // TODO: other fields
    expanded: {},
    sorting: [],
    columnPage: 1,
    rowPage: 1,
    enableComparison: false,
    activeCell: null,
    rowJoinType: "nest",
  };
}
