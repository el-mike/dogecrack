export const VAST_PROVIDER_NAME = 'vast';

export type VastInstanceDto = {
  id: number;
  actual_status: string;
  ssh_host: string;
  ssh_port: string;
  docker_image: string;
  dph_total: number;
  dlperf: number;
  dlperf_per_dphtotal: number;
  gpu_name: string;
  num_gpus: number;
};

export type VastInstance = {
  id: VastInstanceDto['id'];
  status: VastInstanceDto['actual_status'];
  sshHost: VastInstanceDto['ssh_host'];
  sshPort: VastInstanceDto['ssh_port'];
  dockerImage: VastInstanceDto['docker_image'];
  dphTotal: VastInstanceDto['dph_total'];
  dlperfPerDphTotal: VastInstanceDto['dlperf_per_dphtotal'];
  dlperf: VastInstanceDto['dlperf'];
  gpuName: VastInstanceDto['gpu_name'];
  gpuNum: VastInstanceDto['num_gpus'];
};

export const mapVastInstance = (dto: VastInstanceDto) => ({
  id: dto.id,
  status: dto.actual_status,
  sshHost: dto.ssh_host,
  sshPort: dto.ssh_port,
  dockerImage: dto.docker_image,
  dphTotal: dto.dph_total,
  dlperf: dto.dlperf,
  dlperfPerDphTotal: dto.dlperf_per_dphtotal,
  gpuName: dto.gpu_name,
  gpuNum: dto.num_gpus,
} as VastInstance);
