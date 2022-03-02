import { Dictionary } from '../structs/dictionary';

import {
  BaseEntityDto,
  BaseEntity,
} from './base-entity';

import {
  VAST_PROVIDER_NAME,
  VastInstanceDto,
  VastInstance,
  mapVastInstance,
} from './vast-instance';

export type HostInstanceDto = Dictionary & VastInstanceDto;
export type HostInstance = Dictionary & VastInstance;

export type ProgressInfoDto = {
  checked: number;
  total: number;
};

export type ProgressInfo = ProgressInfoDto;

export type PitbullDto = {
  status: number;
  progress: ProgressInfoDto;
  lastOutput: string;
};

export type Pitbull = Omit<PitbullDto, 'progress'> & {
  progress: ProgressInfo;
}

export type PitbullInstanceDto = BaseEntityDto & {
  rules: string[];

  walletString: string;
  passlistUrl: string;

  startedAt: string;
  completedAt: string;

  status: number;

  pitbull: PitbullDto;

  providerName: string;
  hostInstance: HostInstanceDto;
};

export type PitbullInstance = BaseEntity
  & Omit<PitbullInstanceDto, 'pitbull' | 'hostInstance'>
  & {
    pitbull: Pitbull;
    hostInstance: HostInstance;
  };

export const mapPitbullInstance = (dto: PitbullInstanceDto) => ({
  ...dto,
  hostInstance: dto.providerName === VAST_PROVIDER_NAME
    ? mapVastInstance(dto.hostInstance)
    : {},
} as PitbullInstance);
