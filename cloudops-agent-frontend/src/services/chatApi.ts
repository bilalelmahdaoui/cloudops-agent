import type { ApiService } from "./apiService";

interface ChatResponse {
  message: string;
}

export interface ChatHistoryMessage {
  role: "user" | "assistant";
  content: string;
}

export class ChatApi {
  private readonly apiService: ApiService;

  constructor(apiService: ApiService) {
    this.apiService = apiService;
  }

  public async sendMessage(
    message: string,
    history: ChatHistoryMessage[],
  ): Promise<string> {
    const data = await this.apiService.request<ChatResponse>(
      "/chat",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ message, history }),
      },
    );

    return data.message;
  }
}
