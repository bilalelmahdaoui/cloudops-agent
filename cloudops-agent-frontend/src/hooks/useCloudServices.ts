import { useCallback, useEffect, useState } from "react";

import { CloudServiceApi } from "../services/CloudServiceApi";
import type { CloudService } from "../types/CloudService";

export function useCloudServices(api: CloudServiceApi) {
  const [services, setServices] = useState<CloudService[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [restartingServiceId, setRestartingServiceId] =
    useState<string | null>(null);

  const loadServices = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const cloudServices = await api.getCloudServices();

      setServices(cloudServices);
    } catch {
      setError("Impossible de récupérer les services cloud.");
    } finally {
      setLoading(false);
    }
  }, [api]);

  useEffect(() => {
    void loadServices();
  }, [loadServices]);

  const restartService = async (id: string) => {
    try {
      setRestartingServiceId(id);
      setError(null);

      const updatedService = await api.restartCloudService(id);

      setServices((currentServices) =>
        currentServices.map((service) =>
          service.id === updatedService.id
            ? updatedService
            : service,
        ),
      );
    } catch {
      setError("Impossible de redémarrer le service cloud.");
    } finally {
      setRestartingServiceId(null);
    }
  };

  return {
    services,
    loading,
    error,
    restartingServiceId,
    restartService,
  };
}
