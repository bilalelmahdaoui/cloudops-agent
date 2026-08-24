import { useEffect, useState } from "react";

import type { ChatApi } from "../../services/ChatApi";
import type { ChatMessage } from "../../types/ChatMessage";
import { ChatInput } from "./ChatInput";
import { MessageList } from "./MessageList";

const INITIAL_MESSAGES: ChatMessage[] = [
  {
    id: "welcome",
    role: "assistant",
    content:
      "Bonjour. Je peux vous aider à diagnostiquer et piloter vos services cloud.",
  },
];

const CHAT_HISTORY_STORAGE_KEY = "cloudops-agent.chat-history";
const MAX_CONTEXT_MESSAGES = 12;

interface ChatProps {
  api: ChatApi;
  onAgentResponse: () => Promise<void>;
}

export function Chat({
  api,
  onAgentResponse,
}: ChatProps) {
  const [messages, setMessages] =
    useState<ChatMessage[]>(loadChatHistory);
  const [isSending, setIsSending] = useState(false);

  useEffect(() => {
    try {
      localStorage.setItem(
        CHAT_HISTORY_STORAGE_KEY,
        JSON.stringify(messages),
      );
    } catch {
      // Le chat reste fonctionnel si le stockage du navigateur est indisponible.
    }
  }, [messages]);

  const handleSend = async (content: string) => {
    const userMessage: ChatMessage = {
      id: crypto.randomUUID(),
      role: "user",
      content,
    };

    setMessages((currentMessages) => [
      ...currentMessages,
      userMessage,
    ]);

    try {
      setIsSending(true);

      const history = messages
        .slice(-MAX_CONTEXT_MESSAGES)
        .map(({ role, content: messageContent }) => ({
          role,
          content: messageContent,
        }));
      const response = await api.sendMessage(content, history);
      const assistantMessage: ChatMessage = {
        id: crypto.randomUUID(),
        role: "assistant",
        content: response,
      };

      setMessages((currentMessages) => [
        ...currentMessages,
        assistantMessage,
      ]);

      await onAgentResponse();
    } catch {
      const errorMessage: ChatMessage = {
        id: crypto.randomUUID(),
        role: "assistant",
        content:
          "Impossible de contacter CloudOps Agent. Veuillez réessayer.",
      };

      setMessages((currentMessages) => [
        ...currentMessages,
        errorMessage,
      ]);
    } finally {
      setIsSending(false);
    }
  };

  return (
    <section className="chat-panel">
      <header className="panel-header">
        <div>
          <span className="eyebrow">Assistant IA</span>
          <h1>CloudOps Agent</h1>
        </div>

        <span className="agent-status">
          <span className="agent-status__dot" />
          Disponible
        </span>
      </header>

      <MessageList
        messages={messages}
        isThinking={isSending}
      />

      <ChatInput
        onSend={handleSend}
        disabled={isSending}
      />
    </section>
  );
}

function loadChatHistory(): ChatMessage[] {
  try {
    const storedMessages = localStorage.getItem(
      CHAT_HISTORY_STORAGE_KEY,
    );
    if (storedMessages === null) {
      return INITIAL_MESSAGES;
    }

    const parsedMessages: unknown = JSON.parse(storedMessages);
    if (
      !Array.isArray(parsedMessages) ||
      !parsedMessages.every(isChatMessage)
    ) {
      return INITIAL_MESSAGES;
    }

    return parsedMessages.length > 0
      ? parsedMessages
      : INITIAL_MESSAGES;
  } catch {
    return INITIAL_MESSAGES;
  }
}

function isChatMessage(value: unknown): value is ChatMessage {
  if (typeof value !== "object" || value === null) {
    return false;
  }

  const message = value as Record<string, unknown>;
  return (
    typeof message.id === "string" &&
    (message.role === "user" || message.role === "assistant") &&
    typeof message.content === "string"
  );
}
