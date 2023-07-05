import { writable, Writable, Readable, derived } from "svelte/store";
import { getContext } from "svelte";
import type { QueryClient } from "@tanstack/svelte-query";
import type { Runtime } from "@rilldata/web-common/runtime-client/runtime-store";
import {
  MetricsExplorerEntity,
  MetricsExplorerStoreType,
  metricsExplorerStore,
  useDashboardStore,
} from "../dashboard-stores";
import { runtime } from "@rilldata/web-common/runtime-client/runtime-store";

type BusinessModel = {
  runtime: Writable<Runtime>;
  metricsViewName: Writable<string>;
  metricsStore: Readable<MetricsExplorerStoreType>;
  dashboardStore: Readable<MetricsExplorerEntity>;
  queryClient: QueryClient;
  setMetricsViewName: (s: string) => void;
};

export const DEFAULT_STORE_KEY = Symbol("business-model");

export function getBusinessModel(): BusinessModel {
  return getContext(DEFAULT_STORE_KEY);
}

export function createBusinessModel({
  queryClient,
  metricsViewName,
}: {
  queryClient: QueryClient;
  metricsViewName: string;
}): BusinessModel {
  const metricsViewNameStore = writable(metricsViewName);
  const dashboardStore = derived([metricsViewNameStore], ([name], set) => {
    const store = useDashboardStore(name);
    return store.subscribe(set);
  });
  return {
    runtime: runtime,
    metricsViewName: metricsViewNameStore,
    metricsStore: metricsExplorerStore,
    queryClient,
    dashboardStore,
    setMetricsViewName: (name) => {
      metricsViewNameStore.set(name);
    },
  } as BusinessModel;
}
