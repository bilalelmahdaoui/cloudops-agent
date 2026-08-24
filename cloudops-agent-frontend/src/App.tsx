import { useMemo } from "react";

import "./App.css";

import { Chat } from "./components/chat/Chat";
import { CloudServiceList } from "./components/cloud-services/CloudServiceList";
import { useCloudServices } from "./hooks/useCloudServices";
import { ChatApi } from "./services/ChatApi";
import { CloudServiceApi } from "./services/CloudServiceApi";

const API_BASE_URL =
  import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export default function App() {
  const cloudServiceApi = useMemo(
    () => new CloudServiceApi(API_BASE_URL),
    [],
  );
  const chatApi = useMemo(
    () => new ChatApi(API_BASE_URL),
    [],
  );

  const {
    services,
    loading,
    error,
    restartingServiceId,
    restartService,
    refreshServices,
  } = useCloudServices(cloudServiceApi);

  return (
    <main className="app-shell">
      <Chat
        api={chatApi}
        onAgentResponse={refreshServices}
      />

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
