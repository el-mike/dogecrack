export const isString = (candidate: any): candidate is string => {
  return typeof candidate === 'string';
};

export const isNumber = (candidate: any): candidate is number => {
  return typeof candidate === 'number';
};

export const isObject = (candidate: any): candidate is object => {
  return typeof candidate === 'object' && candidate !== null;
};

export const isArray = (candidate: any): candidate is any[] => {
  return Array.isArray(candidate);
};

export const isFunction = (candidate: any): candidate is Function => {
  return typeof candidate === 'function';
};

export const isFile = (candidate: any): candidate is File => {
  return candidate instanceof File;
};

export const isDate = (candidate: any): candidate is Date => {
  return candidate instanceof Date;
};
