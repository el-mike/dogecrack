export type FromKeys<TKeys extends string, TValue> = {
  [K in TKeys]: TValue;
};
