import { useState } from "react";

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

export function Chat() {
  const [messages, setMessages] =
    useState<ChatMessage[]>(INITIAL_MESSAGES);

  const handleSend = (content: string) => {
    const message: ChatMessage = {
      id: crypto.randomUUID(),
      role: "user",
      content,
    };

    setMessages((currentMessages) => [
      ...currentMessages,
      message,
    ]);
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

      <MessageList messages={messages} />

      <ChatInput onSend={handleSend} />
    </section>
  );
}