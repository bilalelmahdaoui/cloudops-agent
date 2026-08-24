import type { CloudService } from "../types/CloudService";

export class CloudServiceApi {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  public async getCloudServices(): Promise<CloudService[]> {
    return this.request<CloudService[]>("/cloud-services");
  }

  public async getCloudService(id: string): Promise<CloudService> {
    return this.request<CloudService>(`/cloud-services/${id}`);
  }

  public async restartCloudService(id: string): Promise<CloudService> {
    return this.request<CloudService>(
      `/cloud-services/${id}/restart`,
      {
        method: "POST",
      },
    );
  }

  public async shutdownCloudService(id: string): Promise<CloudService> {
    return this.request<CloudService>(
      `/cloud-services/${id}/shutdown`,
      {
        method: "POST",
      },
    );
  }

  public async startCloudService(id: string): Promise<CloudService> {
    return this.request<CloudService>(
      `/cloud-services/${id}/start`,
      {
        method: "POST",
      },
    );
  }

  private async request<T>(
    path: string,
    options?: RequestInit,
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}${path}`, options);

    if (!response.ok) {
      throw new Error(
        `The request failed with status ${response.status}`,
      );
    }

    return response.json() as Promise<T>;
  }
}
