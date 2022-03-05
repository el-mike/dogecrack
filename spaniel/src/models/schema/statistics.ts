export type CrackJobsStatisticsDto = {
  all: number;
  acknowledged: number;
  processing: number;
  queued: number;
  rejected: number;
};

export type CrackJobsStatistics = CrackJobsStatisticsDto;


export type PitbullInstancesStatisticsDto = {
  all: number;
  waitingForHost: number;
  hostStarting: number;
  running: number;
  completed: number;
  failed: number;
  interrupted: number;
  success: number;
  passwordsChecked: number;
  totalCost: number;
  averageCost: number;
};

export type PitbullInstancesStatistics = PitbullInstancesStatisticsDto;

export type StatisticsDto = {
  crackJobs: CrackJobsStatisticsDto;
  pitbullInstances: PitbullInstancesStatisticsDto;
}

export type Statistics = {
  crackJobs: CrackJobsStatistics;
  pitbullInstances: PitbullInstancesStatistics;
};
