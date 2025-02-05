import { filterExpressions } from "@rilldata/web-common/features/dashboards/stores/filter-utils";
import {
  ResourceKind,
  useClientFilteredResources,
  useFilteredResources,
  useResource,
} from "@rilldata/web-common/features/entity-management/resource-selectors";
import type {
  RpcStatus,
  V1Expression,
  V1GetResourceResponse,
  V1MetricsViewSpec,
  V1MetricsViewTimeRangeResponse,
  V1Resource,
} from "@rilldata/web-common/runtime-client";
import {
  createQueryServiceMetricsViewTimeRange,
  createRuntimeServiceListResources,
} from "@rilldata/web-common/runtime-client";
import type {
  CreateQueryOptions,
  CreateQueryResult,
} from "@tanstack/svelte-query";
import { derived } from "svelte/store";
import type { ErrorType } from "../../runtime-client/http-client";

export function useMetricsView(
  instanceId: string,
  metricsViewName: string,
  queryOptions?: CreateQueryOptions<
    V1GetResourceResponse,
    ErrorType<RpcStatus>,
    V1Resource
  >,
) {
  return useResource(
    instanceId,
    metricsViewName,
    ResourceKind.MetricsView,
    queryOptions,
  );
}

export function useValidExplores(instanceId: string) {
  // This is used in cloud as well so do not use "useClientFilteredResources"
  return useFilteredResources(instanceId, ResourceKind.Explore, (data) =>
    data?.resources?.filter((res) => !!res.explore?.state?.validSpec),
  );
}

export function useValidCanvases(instanceId: string) {
  return useFilteredResources(instanceId, ResourceKind.Canvas, (data) =>
    data?.resources?.filter((res) => !!res.canvas?.state?.validSpec),
  );
}

export function useValidDashboards(instanceId: string) {
  return createRuntimeServiceListResources(
    instanceId,
    undefined, // TODO: it'd be nice if we could provide multiple kinds here
    {
      query: {
        select: (data) => {
          // Filter for valid Explores and Canvases
          return data?.resources?.filter(
            (res) =>
              !!res.explore?.state?.validSpec || !!res.canvas?.state?.validSpec,
          );
        },
      },
    },
  );
}

/**
 * Gets the valid metrics view spec. Only to be used in displaying a dashboard.
 * Use {@link useDashboard} in the metrics view editor and other use cases.
 */
export const useMetricsViewValidSpec = <T = V1MetricsViewSpec>(
  instanceId: string,
  metricsViewName: string,
  selector?: (meta: V1MetricsViewSpec) => T,
) => {
  return useResource<T>(instanceId, metricsViewName, ResourceKind.MetricsView, {
    select: (data) =>
      selector
        ? selector(data.resource?.metricsView?.state?.validSpec)
        : (data.resource?.metricsView?.state?.validSpec as T),
  });
};

// TODO: cleanup usage of useModelHasTimeSeries and useModelAllTimeRange
export const useModelHasTimeSeries = (
  instanceId: string,
  metricsViewName: string,
) =>
  useMetricsViewValidSpec(
    instanceId,
    metricsViewName,
    (meta) => !!meta?.timeDimension,
  );

export function useMetricsViewTimeRange(
  instanceId: string,
  metricsViewName: string,
  options?: {
    query?: CreateQueryOptions<V1MetricsViewTimeRangeResponse>;
  },
): CreateQueryResult<V1MetricsViewTimeRangeResponse> {
  const { query: queryOptions } = options ?? {};

  return derived(
    [useMetricsViewValidSpec(instanceId, metricsViewName)],
    ([metricsView], set) =>
      createQueryServiceMetricsViewTimeRange(
        instanceId,
        metricsViewName,
        {},
        {
          query: {
            ...queryOptions,
            enabled: !!metricsView.data?.timeDimension && queryOptions?.enabled,
          },
        },
      ).subscribe(set),
  );
}

export const useMetricsViewSpecMeasure = (
  instanceId: string,
  metricsViewName: string,
  measureName: string,
) =>
  useMetricsViewValidSpec(instanceId, metricsViewName, (meta) =>
    meta?.measures?.find((measure) => measure.name === measureName),
  );

/**
 * Returns a copy of a V1MetricsViewFilter that does not include
 * the filters for the specified dimension name.
 */
export const getFiltersForOtherDimensions = (
  filters: V1Expression,
  dimensionName: string,
) => {
  if (!filters) return undefined;
  return filterExpressions(
    filters,
    (e) => e.cond?.exprs?.[0].ident !== dimensionName,
  );
};

export const useGetMetricsViewsForModel = (
  instanceId: string,
  modelName: string,
) => {
  return useClientFilteredResources(
    instanceId,
    ResourceKind.MetricsView,
    (res) =>
      res.metricsView?.spec?.model === modelName ||
      res.metricsView?.spec?.table === modelName,
  );
};

export const useGetExploresForMetricsView = (
  instanceId: string,
  metricsViewName: string,
) => {
  return useClientFilteredResources(
    instanceId,
    ResourceKind.Explore,
    (res) => res.explore?.spec?.metricsView === metricsViewName,
  );
};
