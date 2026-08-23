import { useMemo } from "react";

import "./App.css";

import { Chat } from "./components/chat/Chat";
import { CloudServiceList } from "./components/cloud-services/CloudServiceList";
import { useCloudServices } from "./hooks/useCloudServices";
import { CloudServiceApi } from "./services/CloudServiceApi";

export default function App() {
  const cloudServiceApi = useMemo(
    () =>
      new CloudServiceApi(
        import.meta.env.VITE_API_URL ??
          "http://localhost:8080",
      ),
    [],
  );

  const {
    services,
    loading,
    error,
    restartingServiceId,
    restartService,
  } = useCloudServices(cloudServiceApi);

  return (
    <main className="app-shell">
      <Chat />

      <aside className="services-panel">
        <header className="panel-header">
          <div>
            <span className="eyebrow">Infrastructure</span>
            <h2>Services cloud</h2>
          </div>

          <span className="service-count">
            {services.length}
          </span>
        </header>

        <CloudServiceList
          services={services}
          loading={loading}
          error={error}
          restartingServiceId={restartingServiceId}
          onRestart={restartService}
        />
      </aside>
    </main>
  );
}