export type PitbullInstancesStatisticsDto = {
  all: number;
  waitingForHost: number;
  hostStarting: number;
  running: number;
  completed: number;
  failed: number;
  interrupted: number;
  success: number;
};

export type PitbullInstancesStatistics = PitbullInstancesStatisticsDto;
