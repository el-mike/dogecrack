import {
  BaseEntityDto,
  BaseEntity,
} from './base-entity';

import {
  PitbullInstanceDto,
  PitbullInstance,
  mapPitbullInstance,
} from './pitbull-instance';

export type PitbullJobDto = BaseEntityDto & {
  keyword: string;
  walletString: string;
  passlistUrl: string;

  status: number;

  instanceId: string;
  instance: PitbullInstanceDto;

  startedAt: string;
  firstScheduledAt: string;
  lastScheduledAt: string;
  acknowledgedAt: string;
  rejectedAt: string;

  rescheduleCount: number;

  errorLog: string;
};

export type PitbullJob = BaseEntity
  & Omit<PitbullJobDto, 'instance'>
  & {
    instance: PitbullInstance;
  };

export const mapPitbullJob = (dto: PitbullJobDto) => ({
  ...dto,
  instance: mapPitbullInstance(dto.instance),
} as PitbullJob);
