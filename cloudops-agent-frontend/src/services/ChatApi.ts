interface ChatResponse {
  message: string;
}

export interface ChatHistoryMessage {
  role: "user" | "assistant";
  content: string;
}

export class ChatApi {
  private readonly baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  public async sendMessage(
    message: string,
    history: ChatHistoryMessage[],
  ): Promise<string> {
    const response = await fetch(`${this.baseUrl}/chat`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ message, history }),
    });

    if (!response.ok) {
      throw new Error(
        `The request failed with status ${response.status}`,
      );
    }

    const data = (await response.json()) as ChatResponse;
    return data.message;
  }
}
