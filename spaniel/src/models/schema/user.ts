import { BaseEntity } from './base-entity';

export type User = BaseEntity & {
  name: string;
  email: string;
  role: string;
};
