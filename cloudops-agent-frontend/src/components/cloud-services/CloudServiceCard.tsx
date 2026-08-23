import { useState } from "react";

import {
  CloudServiceStatus,
  type CloudService,
} from "../../types/CloudService";

interface CloudServiceCardProps {
  service: CloudService;
  isRestarting: boolean;
  onRestart: (id: string) => Promise<void>;
}

export function CloudServiceCard({
  service,
  isRestarting,
  onRestart,
}: CloudServiceCardProps) {
  const [isConfirming, setIsConfirming] = useState(false);

  const cpuPercentage = Math.round(service.cpuUsage * 100);

  const handleRestart = async () => {
    setIsConfirming(false);
    await onRestart(service.id);
  };

  return (
    <article className="service-card">
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
        <span>{service.logs.length} événements</span>

        {!isConfirming ? (
          <button
            type="button"
            className="button button--secondary"
            disabled={isRestarting}
            onClick={() => setIsConfirming(true)}
          >
            {isRestarting ? "Redémarrage..." : "Redémarrer"}
          </button>
        ) : (
          <div className="confirmation-actions">
            <button
              type="button"
              className="button button--danger"
              onClick={() => void handleRestart()}
            >
              Confirmer
            </button>

            <button
              type="button"
              className="button button--ghost"
              onClick={() => setIsConfirming(false)}
            >
              Annuler
            </button>
          </div>
        )}
      </div>
    </article>
  );
}

function formatStatus(status: CloudServiceStatus): string {
  switch (status) {
    case CloudServiceStatus.Running:
      return "En ligne";

    case CloudServiceStatus.Restarting:
      return "Redémarrage";

    case CloudServiceStatus.Down:
      return "Indisponible";
  }
}