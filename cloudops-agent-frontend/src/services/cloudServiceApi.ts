import type { CloudService } from "../types/cloudService";

export class CloudServiceApi {
  constructor(private readonly baseUrl: string) {}

  public async getCloudService(id: string): Promise<CloudService> {
    const response = await fetch(`${this.baseUrl}/cloud-services/${id}`);

    if (!response.ok) {
      throw new Error("Impossible de récupérer le service cloud");
    }

    return response.json();
  }

  public async restartCloudService(id: string): Promise<CloudService> {
    const response = await fetch(
      `${this.baseUrl}/cloud-services/${id}/restart`,
      {
        method: "POST",
      }
    );

    if (!response.ok) {
      throw new Error("Impossible de redémarrer le service cloud");
    }

    return response.json();
  }
}