import type { DimensionDefinitionEntity } from "$common/data-modeler-state-service/entity-state-service/DimensionDefinitionStateService";
import {
  EntityStatus,
  EntityType,
} from "$common/data-modeler-state-service/entity-state-service/EntityStateService";
import type { MeasureDefinitionEntity } from "$common/data-modeler-state-service/entity-state-service/MeasureDefinitionStateService";
import type { TimeSeriesTimeRange } from "$common/database-service/DatabaseTimeSeriesActions";
import { getArrayDiff } from "$common/utils/getArrayDiff";
import { generateBigNumbersApi } from "$lib/redux-store/big-number/big-number-apis";
import {
  addDimensionToExplore,
  addMeasureToExplore,
  clearSelectedLeaderboardValues,
  initMetricsExplorer,
  LeaderboardValues,
  MetricsExplorerEntity,
  removeDimensionFromExplore,
  removeMeasureFromExplore,
  setExplorerIsStale,
  setExplorerSelectableTimeRange,
  setExploreSelectedTimeRange,
  setExploreTimeRange,
  setLeaderboardDimensionValues,
  setLeaderboardMeasureId,
  setLeaderboardValuesErrorStatus,
  setLeaderboardValuesStatus,
  toggleExploreMeasure,
  toggleLeaderboardActiveValue,
} from "$lib/redux-store/explore/explore-slice";
import { selectValidMeasures } from "$lib/redux-store/measure-definition/measure-definition-selectors";
import { createAsyncThunk } from "$lib/redux-store/redux-toolkit-wrapper";
import type { RillReduxState } from "$lib/redux-store/store-root";
import { generateTimeSeriesApi } from "$lib/redux-store/timeseries/timeseries-apis";
import { fetchWrapper, streamingFetchWrapper } from "$lib/util/fetchWrapper";
import {
  getDefaultSelectedTimeRange,
  makeSelectableTimeRanges,
  prune,
} from "../../../routes/_surfaces/workspace/explore/utils";
import { selectMetricsExplorerById } from "$lib/redux-store/explore/explore-selectors";
import { setReferenceValues } from "$lib/redux-store/big-number/big-number-slice";

/**
 * A wrapper to dispatch updates to explore.
 * Currently, it updates these sections
 * 1. Leaderboard values based on selected filters and measure
 * 2. Time series for all selected measures
 * 3. Big numbers for all selected measures
 */
const updateExploreWrapper = (dispatch, metricsDefId: string) => {
  dispatch(generateTimeSeriesApi({ id: metricsDefId }));
  dispatch(generateBigNumbersApi({ id: metricsDefId }));
  dispatch(updateLeaderboardValuesApi(metricsDefId));
};

/**
 * Syncs explore with updated measures and dimensions.
 * If a MetricsExplorer entity is not present then a new one is created.
 * It then calls {@link updateExploreWrapper} to update explore.
 * It also dispatches {@link fetchTimestampColumnRangeApi} to update time range.
 */
export const syncExplore = async (
  dispatch,
  metricsDefId: string,
  metricsExplorer: MetricsExplorerEntity,
  dimensions: Array<DimensionDefinitionEntity>,
  measures: Array<MeasureDefinitionEntity>
) => {
  if (measures) measures = selectValidMeasures(measures);

  let shouldUpdate = false;
  if (!metricsExplorer) {
    dispatch(initMetricsExplorer(metricsDefId, dimensions, measures));
    shouldUpdate = true;
  } else {
    if (dimensions)
      shouldUpdate = syncDimensions(dispatch, metricsExplorer, dimensions);
    if (measures)
      shouldUpdate ||= syncMeasures(dispatch, metricsExplorer, measures);
  }

  // To avoid infinite loop only update if something changed.
  if (shouldUpdate || metricsExplorer.isStale) {
    dispatch(setExplorerIsStale(metricsDefId, false));
    await dispatch(fetchTimestampColumnRangeApi(metricsDefId));
    updateExploreWrapper(dispatch, metricsDefId);
  }
};
/**
 * Syncs dimensions from MetricsDefinition with MetricsExplorer dimensions.
 * Calls {@link addDimensionToExplore} for missing dimension.
 * Calls {@link removeDimensionFromExplore} for excess dimension.
 */
const syncDimensions = (
  dispatch,
  metricsExplorer: MetricsExplorerEntity,
  dimensions: Array<DimensionDefinitionEntity>
) => {
  const { extraSrc: addDimensions, extraTarget: removeDimensions } =
    getArrayDiff(
      dimensions,
      (dimension) => dimension.id,
      metricsExplorer.leaderboards,
      (leaderboard) => leaderboard.dimensionId
    );
  addDimensions.forEach((addDimension) =>
    dispatch(addDimensionToExplore(metricsExplorer.id, addDimension.id))
  );
  removeDimensions.forEach((removeDimension) =>
    dispatch(
      removeDimensionFromExplore(
        metricsExplorer.id,
        removeDimension.dimensionId
      )
    )
  );

  return addDimensions.length > 0 || removeDimensions.length > 0;
};
/**
 * Syncs measures from MetricsDefinition with MetricsExplorer measures.
 * Calls {@link addMeasureToExplore} for missing measure.
 * Calls {@link removeMeasureFromExplore} for excess measure.
 */
const syncMeasures = (
  dispatch,
  metricsExplorer: MetricsExplorerEntity,
  measures: Array<MeasureDefinitionEntity>
) => {
  const { extraSrc: addMeasures, extraTarget: removeMeasures } = getArrayDiff(
    measures,
    (measure) => measure.id,
    metricsExplorer.measureIds,
    (measureId) => measureId
  );
  addMeasures.forEach((addMeasure) =>
    dispatch(addMeasureToExplore(metricsExplorer.id, addMeasure.id))
  );
  removeMeasures.forEach((removeMeasure) =>
    dispatch(removeMeasureFromExplore(metricsExplorer.id, removeMeasure))
  );

  return addMeasures.length > 0 || removeMeasures.length > 0;
};

/**
 * Toggles selection of a measures to be displayed.
 * It then updates,
 * 1. Time series for all selected measures
 * 2. Big numbers for all selected measures
 */
export const toggleExploreMeasureAndUpdate = (
  dispatch,
  metricsDefId: string,
  selectedMeasureId: string
) => {
  dispatch(toggleExploreMeasure(metricsDefId, selectedMeasureId));
  dispatch(generateTimeSeriesApi({ id: metricsDefId }));
  dispatch(generateBigNumbersApi({ id: metricsDefId }));
};

/**
 * Sets the measure id used in leaderboard for ranking and other calculations.
 * It then updates Leaderboard values based on selected filters and measure
 */
export const setMeasureIdAndUpdateLeaderboard = (
  dispatch,
  metricsDefId: string,
  measureId: string
) => {
  dispatch(setLeaderboardMeasureId(metricsDefId, measureId));
  dispatch(updateLeaderboardValuesApi(metricsDefId));
};

/**
 * Toggles a selected value in the leaderboard.
 * Pass 'include' param boolean to denote whether the value is included or excluded in time series and big number queries.
 * It then calls {@link updateExploreWrapper} to update explore.
 */
export const toggleSelectedLeaderboardValueAndUpdate = (
  dispatch,
  metricsDefId: string,
  dimensionId: string,
  dimensionValue: unknown,
  include: boolean
) => {
  dispatch(
    toggleLeaderboardActiveValue(
      metricsDefId,
      dimensionId,
      dimensionValue,
      include
    )
  );
  updateExploreWrapper(dispatch, metricsDefId);
};

/**
 * Clears all selected values in the leaderboard.
 * It then calls {@link updateExploreWrapper} to update explore.
 */
export const clearSelectedLeaderboardValuesAndUpdate = (
  dispatch,
  metricsDefId: string
) => {
  dispatch(clearSelectedLeaderboardValues(metricsDefId));
  updateExploreWrapper(dispatch, metricsDefId);
};

/**
 * Sets user selected time rage.
 * It then calls {@link updateExploreWrapper} to update explore.
 */
export const setExploreSelectedTimeRangeAndUpdate = (
  dispatch,
  metricsDefId: string,
  selectedTimeRange: Partial<TimeSeriesTimeRange>
) => {
  dispatch(setReferenceValues(metricsDefId, undefined));
  dispatch(setExploreSelectedTimeRange(metricsDefId, selectedTimeRange));
  updateExploreWrapper(dispatch, metricsDefId);
};

/**
 * Async-thunk to update leaderboard values.
 * Streams dimension values from backend per dimension and updates it in the state.
 */
export const updateLeaderboardValuesApi = createAsyncThunk(
  `${EntityType.MetricsExplorer}/updateLeaderboard`,
  async (metricsDefId: string, thunkAPI) => {
    const state = thunkAPI.getState() as RillReduxState;
    const metricsExplorer: MetricsExplorerEntity = (
      thunkAPI.getState() as RillReduxState
    ).metricsExplorer.entities[metricsDefId];
    const filters = prune(
      metricsExplorer.activeValues,
      state.dimensionDefinition.entities
    );
    const requestBody = {
      measureId: metricsExplorer.leaderboardMeasureId,
      filters,
      timeRange: metricsExplorer.selectedTimeRange,
    };

    thunkAPI.dispatch(
      setLeaderboardValuesStatus(metricsDefId, EntityStatus.Running)
    );

    const stream = streamingFetchWrapper<LeaderboardValues>(
      `metrics/${metricsExplorer.id}/leaderboards`,
      "POST",
      requestBody
    );
    for await (const dimensionData of stream) {
      thunkAPI.dispatch(
        setLeaderboardDimensionValues(
          metricsExplorer.id,
          dimensionData.dimensionId,
          dimensionData.values
        )
      );
    }

    thunkAPI.dispatch(setLeaderboardValuesErrorStatus(metricsDefId));
  }
);

/**
 * Fetches time range for the selected timestamp column.
 * Store the response in MetricsExplorer slice by calling {@link setExploreTimeRange}
 */
export const fetchTimestampColumnRangeApi = createAsyncThunk(
  `${EntityType.MetricsExplorer}/getTimestampColumnRange`,
  async (metricsDefId: string, thunkAPI) => {
    const timeRange = await fetchWrapper(
      `metrics/${metricsDefId}/time-range`,
      "GET"
    );
    thunkAPI.dispatch(setExploreTimeRange(metricsDefId, timeRange));
    const selectableTimeRanges = makeSelectableTimeRanges(timeRange);
    thunkAPI.dispatch(
      setExplorerSelectableTimeRange(metricsDefId, selectableTimeRanges)
    );

    const metricsExplorer = selectMetricsExplorerById(
      thunkAPI.getState() as RillReduxState,
      metricsDefId
    );
    if (!metricsExplorer.selectedTimeRange) {
      const selectedTimeRange =
        getDefaultSelectedTimeRange(selectableTimeRanges);
      setExploreSelectedTimeRangeAndUpdate(thunkAPI.dispatch, metricsDefId, {
        name: selectedTimeRange.name,
        start: selectedTimeRange.start,
        end: selectedTimeRange.end,
      });
    }
  }
);
