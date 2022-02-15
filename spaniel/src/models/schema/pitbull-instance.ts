import { Dictionary } from '../structs/dictionary';
import { BaseEntity } from './base-entity';

export type HostInstance = Dictionary;

export type ProgressInfo = {
  checked: number;
  total: number;
};

export type PitbullInstance = BaseEntity & {
  rules: string[];

  walletString: string;
  passlistUrl: string;

  status: number;
  progress: ProgressInfo;
  lastOutput: string;

  providerName: string;
  hostInstance: HostInstance;
};
