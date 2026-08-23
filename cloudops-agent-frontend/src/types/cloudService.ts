export enum CloudServiceStatus {
  Running = "running",
  Restarting = "restarting",
  Down = "down",
}

export interface CloudService {
  id: string
  name: string
  status: CloudServiceStatus
  cpuUsage: number
  logs: CloudServiceLog[]
}

export interface CloudServiceLog {
  dateTime: string
  event: string
}