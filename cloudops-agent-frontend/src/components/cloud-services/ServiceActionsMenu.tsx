import { useEffect, useRef, useState } from "react";

import { CloudServiceStatus } from "../../types/cloudService";

export type ServiceAction = "restart" | "shutdown" | "start";

interface ServiceActionsMenuProps {
  serviceName: string;
  status: CloudServiceStatus;
  disabled: boolean;
  onAction: (action: ServiceAction) => Promise<void>;
}

export function ServiceActionsMenu({
  serviceName,
  status,
  disabled,
  onAction,
}: ServiceActionsMenuProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [pendingAction, setPendingAction] =
    useState<ServiceAction | null>(null);

  useEffect(() => {
    if (!isOpen) {
      return;
    }

    const closeOnOutsideClick = (event: PointerEvent) => {
      if (
        event.target instanceof Node &&
        !containerRef.current?.contains(event.target)
      ) {
        setIsOpen(false);
        setPendingAction(null);
      }
    };
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setIsOpen(false);
        setPendingAction(null);
      }
    };

    document.addEventListener("pointerdown", closeOnOutsideClick);
    document.addEventListener("keydown", closeOnEscape);
    return () => {
      document.removeEventListener("pointerdown", closeOnOutsideClick);
      document.removeEventListener("keydown", closeOnEscape);
    };
  }, [isOpen]);

  const handleAction = async () => {
    if (pendingAction === null) {
      return;
    }

    const action = pendingAction;
    setIsOpen(false);
    setPendingAction(null);
    await onAction(action);
  };

  return (
    <div ref={containerRef} className="service-actions">
      <button
        type="button"
        className="service-actions__trigger"
        aria-label={`Actions for ${serviceName}`}
        aria-haspopup="menu"
        aria-expanded={isOpen}
        disabled={disabled}
        onClick={() => {
          setIsOpen((current) => !current);
          setPendingAction(null);
        }}
      >
        <span aria-hidden="true">⋯</span>
      </button>

      {isOpen && !disabled && (
        <div
          className="service-actions__menu"
          role={pendingAction === null ? "menu" : "dialog"}
          aria-label={
            pendingAction === null
              ? `Actions for ${serviceName}`
              : confirmationLabel(pendingAction)
          }
        >
          {pendingAction === null ? (
            <ActionChoices
              status={status}
              onSelect={setPendingAction}
            />
          ) : (
            <div className="service-actions__confirmation">
              <p>{confirmationLabel(pendingAction)}</p>
              <div>
                <button
                  type="button"
                  className="service-actions__confirm"
                  onClick={() => void handleAction()}
                >
                  Confirm
                </button>
                <button
                  type="button"
                  onClick={() => setPendingAction(null)}
                >
                  Cancel
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

interface ActionChoicesProps {
  status: CloudServiceStatus;
  onSelect: (action: ServiceAction) => void;
}

function ActionChoices({ status, onSelect }: ActionChoicesProps) {
  if (status === CloudServiceStatus.Down) {
    return (
      <button
        type="button"
        role="menuitem"
        onClick={() => onSelect("start")}
      >
        Start
      </button>
    );
  }

  return (
    <>
      <button
        type="button"
        role="menuitem"
        onClick={() => onSelect("restart")}
      >
        Restart
      </button>
      <button
        type="button"
        role="menuitem"
        className="service-actions__danger"
        onClick={() => onSelect("shutdown")}
      >
        Stop
      </button>
    </>
  );
}

function confirmationLabel(action: ServiceAction): string {
  switch (action) {
    case "restart":
      return "Confirm restart?";
    case "shutdown":
      return "Confirm shutdown?";
    case "start":
      return "Confirm start?";
  }
}
