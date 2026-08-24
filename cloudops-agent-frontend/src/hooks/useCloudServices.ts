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
  const [operatingServiceId, setOperatingServiceId] =
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
      setError("Unable to load cloud services.");
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

  const updateService = (updatedService: CloudService) => {
    setServices((currentServices) =>
      currentServices.map((service) =>
        service.id === updatedService.id
          ? updatedService
          : service,
      ),
    );
  };

  const restartService = async (id: string) => {
    try {
      setOperatingServiceId(id);
      setError(null);

      const updatedService = await api.restartCloudService(id);
      updateService(updatedService);
    } catch {
      setError("Unable to restart the cloud service.");
    } finally {
      setOperatingServiceId(null);
    }
  };

  const shutdownService = async (id: string) => {
    try {
      setOperatingServiceId(id);
      setError(null);

      const updatedService = await api.shutdownCloudService(id);
      updateService(updatedService);
    } catch {
      setError("Unable to stop the cloud service.");
    } finally {
      setOperatingServiceId(null);
    }
  };

  const startService = async (id: string) => {
    try {
      setOperatingServiceId(id);
      setError(null);

      const updatedService = await api.startCloudService(id);
      updateService(updatedService);
    } catch {
      setError("Unable to start the cloud service.");
    } finally {
      setOperatingServiceId(null);
    }
  };

  return {
    services,
    loading,
    error,
    operatingServiceId,
    restartService,
    shutdownService,
    startService,
    refreshServices,
  };
}
