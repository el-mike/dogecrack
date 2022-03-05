/**
 * If number has decimal part, round it to 2 digits.
 * Otherwise, return the number.
 */
export const roundDecimals = (target: number) =>
  (target % 1 !== 0)
    ? target.toFixed(2)
    : target;
