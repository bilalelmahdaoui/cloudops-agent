import {
  CloudServiceStatus,
  type CloudService,
} from "../../types/CloudService";
import { CloudServiceCard } from "./CloudServiceCard";

interface CloudServiceListProps {
  services: CloudService[];
  loading: boolean;
  error: string | null;
  operatingServiceId: string | null;
  onRestart: (id: string) => Promise<void>;
  onShutdown: (id: string) => Promise<void>;
  onStart: (id: string) => Promise<void>;
}

export function CloudServiceList({
  services,
  loading,
  error,
  operatingServiceId,
  onRestart,
  onShutdown,
  onStart,
}: CloudServiceListProps) {
  if (loading) {
    return <p className="panel-state">Loading services...</p>;
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
          isOperating={
            operatingServiceId === service.id ||
            service.status === CloudServiceStatus.Restarting
          }
          onRestart={onRestart}
          onShutdown={onShutdown}
          onStart={onStart}
        />
      ))}
    </div>
  );
}
