import { Dictionary } from '../structs';

import {
  BaseEntityDto,
  BaseEntity,
} from './base-entity';

export type ListRequest<TFilters extends Dictionary = Dictionary> = {
  [K in keyof TFilters]: TFilters[K];
} & {
  page?: number;
  pageSize?: number;
};

export type ListResponse<TEntityDto extends BaseEntityDto = BaseEntityDto> = {
  data: TEntityDto[];
  page: number;
  totalCount: number;
};

export type ListResponseResult<TEntity extends BaseEntity = BaseEntity> =
  Omit<ListResponse, 'data'>
  & {
    entities: TEntity[];
  };
