import {
  CloudServiceStatus,
  type CloudService,
} from "../../types/CloudService";
import { CloudServiceCard } from "./CloudServiceCard";

interface CloudServiceListProps {
  services: CloudService[];
  loading: boolean;
  error: string | null;
  restartingServiceId: string | null;
  onRestart: (id: string) => Promise<void>;
}

export function CloudServiceList({
  services,
  loading,
  error,
  restartingServiceId,
  onRestart,
}: CloudServiceListProps) {
  if (loading) {
    return <p className="panel-state">Chargement des services...</p>;
  }

  if (error) {
    return <p className="panel-state panel-state--error">{error}</p>;
  }

  return (
    <div className="service-list">
      {services.map((service) => (
        <CloudServiceCard
          key={service.id}
          service={service}
          isRestarting={
            restartingServiceId === service.id ||
            service.status === CloudServiceStatus.Restarting
          }
          onRestart={onRestart}
        />
      ))}
    </div>
  );
}
