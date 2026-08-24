import { useEffect, useRef } from "react";

import type { ChatMessage } from "../../types/ChatMessage";
import { MessageBubble } from "./MessageBubble";
import { ThinkingBubble } from "./ThinkingBubble";

interface MessageListProps {
  messages: ChatMessage[];
  isThinking: boolean;
}

export function MessageList({
  messages,
  isThinking,
}: MessageListProps) {
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const list = listRef.current;
    if (list === null) {
      return;
    }

    list.scrollTo({
      top: list.scrollHeight,
      behavior: "smooth",
    });
  }, [messages, isThinking]);

  return (
    <div
      ref={listRef}
      className="message-list"
      aria-live="polite"
    >
      {messages.map((message) => (
        <MessageBubble
          key={message.id}
          message={message}
        />
      ))}

      {isThinking && <ThinkingBubble />}
    </div>
  );
}
