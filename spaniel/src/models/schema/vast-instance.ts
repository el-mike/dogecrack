export const VAST_PROVIDER_NAME = 'vast';

export type VastInstanceDto = {
  id: number;
  actual_status: string;
  ssh_host: string;
  ssh_port: string;
  docker_image: string;
};

export type VastInstance = {
  id: VastInstanceDto['id'];
  status: VastInstanceDto['actual_status'];
  sshHost: VastInstanceDto['ssh_host'];
  sshPort: VastInstanceDto['ssh_port'];
  dockerImage: VastInstanceDto['docker_image'];
};

export const mapVastInstance = (dto: VastInstanceDto) => ({
  id: dto.id,
  status: dto.actual_status,
  sshHost: dto.ssh_host,
  sshPort: dto.ssh_port,
  dockerImage: dto.docker_image,
} as VastInstance);
