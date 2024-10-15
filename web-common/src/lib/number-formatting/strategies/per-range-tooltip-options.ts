import {
  type FormatterOptionsCommon,
  type FormatterRangeSpecsStrategy,
  NumberKind,
  type RangeFormatSpec,
} from "../humanizer-types";

const tooltipRangeSpec: RangeFormatSpec[] = [
  {
    minMag: -4,
    supMag: -2,
    maxDigitsRight: 2,
    baseMagnitude: 0,
    overrideValue: {
      int: "",
      dot: ".",
      frac: "00",
      prefix: "~",
      suffix: "",
    },
  },
  {
    minMag: -2,
    supMag: 3,
    maxDigitsRight: 2,
    baseMagnitude: 0,
    useTrailingDot: false,
    padWithInsignificantZeros: false,
  },
  {
    minMag: 3,
    supMag: 11,
    maxDigitsRight: 0,
    maxDigitsLeft: 12,
    baseMagnitude: 0,
    useTrailingDot: false,
    padWithInsignificantZeros: false,
  },
];

export const tooltipNoFormattingOptions: FormatterOptionsCommon &
  FormatterRangeSpecsStrategy = {
  numberKind: NumberKind.ANY,
  rangeSpecs: tooltipRangeSpec,
  defaultMaxDigitsRight: 2,
  upperCaseEForExponent: true,
};

export const tooltipPercentOptions: FormatterOptionsCommon &
  FormatterRangeSpecsStrategy = {
  rangeSpecs: tooltipRangeSpec,
  defaultMaxDigitsRight: 2,
  numberKind: NumberKind.PERCENT,
};

export const tooltipCurrencyOptions = (
  numberKind: NumberKind,
): FormatterOptionsCommon & FormatterRangeSpecsStrategy => ({
  rangeSpecs: tooltipRangeSpec,
  defaultMaxDigitsRight: 2,
  numberKind,
});
