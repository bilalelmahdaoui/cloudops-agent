import {
  CloudServiceStatus,
  type CloudService,
} from "../../types/cloudService";
import {
  ServiceActionsMenu,
  type ServiceAction,
} from "./ServiceActionsMenu";

interface CloudServiceCardProps {
  service: CloudService;
  isOperating: boolean;
  onRestart: (id: string) => Promise<void>;
  onShutdown: (id: string) => Promise<void>;
  onStart: (id: string) => Promise<void>;
}

export function CloudServiceCard({
  service,
  isOperating,
  onRestart,
  onShutdown,
  onStart,
}: CloudServiceCardProps) {
  const isRestarting =
    service.status === CloudServiceStatus.Restarting;
  const isMuted =
    isOperating ||
    isRestarting ||
    service.status === CloudServiceStatus.Down;
  const cpuPercentage =
    service.status === CloudServiceStatus.Running && !isOperating
      ? Math.round(service.cpuUsage * 100)
      : 0;

  const handleAction = (action: ServiceAction) => {
    switch (action) {
      case "restart":
        return onRestart(service.id);
      case "shutdown":
        return onShutdown(service.id);
      case "start":
        return onStart(service.id);
    }
  };

  return (
    <article
      className={`service-card${
        isMuted ? " service-card--muted" : ""
      }`}
    >
      <header className="service-card__header">
        <div>
          <h3>{service.name}</h3>
          <span className="service-card__id">
            {service.id}
          </span>
        </div>

        <span
          className={`status status--${service.status}`}
        >
          <span className="status__dot" />
          {formatStatus(service.status)}
        </span>
      </header>

      <div className="service-card__metric">
        <div className="service-card__metric-header">
          <span>CPU</span>
          <strong>{cpuPercentage}%</strong>
        </div>

        <div className="cpu-bar">
          <div
            className="cpu-bar__value"
            style={{
              width: `${Math.min(cpuPercentage, 100)}%`,
            }}
          />
        </div>
      </div>

      <div className="service-card__footer">
        <span>{service.logs.length} events</span>

        <ServiceActionsMenu
          serviceName={service.name}
          status={service.status}
          disabled={isOperating}
          onAction={handleAction}
        />
      </div>
    </article>
  );
}

function formatStatus(status: CloudServiceStatus): string {
  switch (status) {
    case CloudServiceStatus.Running:
      return "Online";

    case CloudServiceStatus.Restarting:
      return "Restarting";

    case CloudServiceStatus.Down:
      return "Offline";
  }
}
