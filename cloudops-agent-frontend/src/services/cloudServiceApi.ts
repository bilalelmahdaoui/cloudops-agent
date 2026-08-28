import type { CloudService } from "../types/cloudService";
import type { ApiService } from "./apiService";

export class CloudServiceApi {
  private readonly apiService: ApiService;

  constructor(apiService: ApiService) {
    this.apiService = apiService;
  }

  public async getCloudServices(): Promise<CloudService[]> {
    return this.apiService.request<CloudService[]>("/cloud-services");
  }

  public async getCloudService(id: string): Promise<CloudService> {
    return this.apiService.request<CloudService>(`/cloud-services/${id}`);
  }

  public async restartCloudService(id: string): Promise<CloudService> {
    return this.apiService.request<CloudService>(
      `/cloud-services/${id}/restart`,
      {
        method: "POST",
      },
    );
  }

  public async shutdownCloudService(id: string): Promise<CloudService> {
    return this.apiService.request<CloudService>(
      `/cloud-services/${id}/shutdown`,
      {
        method: "POST",
      },
    );
  }

  public async startCloudService(id: string): Promise<CloudService> {
    return this.apiService.request<CloudService>(
      `/cloud-services/${id}/start`,
      {
        method: "POST",
      },
    );
  }
}
