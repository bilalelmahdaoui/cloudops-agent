import { useMemo } from "react";

import "./App.css";

import { Chat } from "./components/chat/Chat";
import { CloudServiceList } from "./components/cloud-services/CloudServiceList";
import { useCloudServices } from "./hooks/useCloudServices";
import { ChatApi } from "./services/ChatApi";
import { CloudServiceApi } from "./services/CloudServiceApi";
import { config } from "./config";

export default function App() {
  const cloudServiceApi = useMemo(
    () => new CloudServiceApi(config.apiBaseUrl),
    [],
  );
  const chatApi = useMemo(
    () => new ChatApi(config.apiBaseUrl),
    [],
  );

  const {
    services,
    loading,
    error,
    operatingServiceId,
    restartService,
    shutdownService,
    startService,
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
            <h2>Cloud services</h2>
          </div>

          <span className="service-count">
            {services.length}
          </span>
        </header>

        <CloudServiceList
          services={services}
          loading={loading}
          error={error}
          operatingServiceId={operatingServiceId}
          onRestart={restartService}
          onShutdown={shutdownService}
          onStart={startService}
        />
      </aside>
    </main>
  );
}
