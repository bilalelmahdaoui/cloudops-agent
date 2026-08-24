import {
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";

import { CloudServiceApi } from "../services/CloudServiceApi";
import type { CloudService } from "../types/CloudService";

const SERVICES_REFRESH_INTERVAL_MS = 1_000;

export function useCloudServices(api: CloudServiceApi) {
  const [services, setServices] = useState<CloudService[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [restartingServiceId, setRestartingServiceId] =
    useState<string | null>(null);
  const refreshInProgress = useRef(false);

  const refreshServices = useCallback(async () => {
    if (refreshInProgress.current) {
      return;
    }

    refreshInProgress.current = true;
    try {
      const cloudServices = await api.getCloudServices();

      setServices(cloudServices);
      setError(null);
    } catch {
      setError("Impossible de récupérer les services cloud.");
    } finally {
      setLoading(false);
      refreshInProgress.current = false;
    }
  }, [api]);

  useEffect(() => {
    const initialLoadId = window.setTimeout(
      () => void refreshServices(),
      0,
    );
    const intervalId = window.setInterval(
      () => void refreshServices(),
      SERVICES_REFRESH_INTERVAL_MS,
    );

    return () => {
      window.clearTimeout(initialLoadId);
      window.clearInterval(intervalId);
    };
  }, [refreshServices]);

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
    refreshServices,
  };
}
