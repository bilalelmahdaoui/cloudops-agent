export const CloudServiceStatus = {
  Running: "running",
  Restarting: "restarting",
  Down: "down",
} as const;

export type CloudServiceStatus =
  (typeof CloudServiceStatus)[keyof typeof CloudServiceStatus];

export interface CloudServiceLog {
  dateTime: string;
  event: string;
}

export interface CloudService {
  id: string;
  name: string;
  status: CloudServiceStatus;
  cpuUsage: number;
  logs: CloudServiceLog[];
}
