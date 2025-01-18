import {
  BaseEntityDto,
  BaseEntity,
} from './base-entity';

export type SettingsDto = BaseEntityDto & {
  startHostAttemptsLimit: number;
  checkStatusRetryLimit: number;
  stalledProgressLimit: number;
  rescheduleLimit: number;
  runningInstancesLimit: number;
  checkHostInterval: number;
  checkPitbullInterval: number;
  vastSearchCriteria: string;
  minPasswordLength: number;
  maxPasswordLength: number;
};

export type Settings = BaseEntity & SettingsDto;
