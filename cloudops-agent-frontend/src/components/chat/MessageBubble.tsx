import type { ChatMessage } from "../../types/ChatMessage";

interface MessageBubbleProps {
  message: ChatMessage;
}

export function MessageBubble({
  message,
}: MessageBubbleProps) {
  return (
    <div
      className={`message message--${message.role}`}
    >
      <span className="message__author">
        {message.role === "user" ? "Vous" : "CloudOps Agent"}
      </span>

      <p>{message.content}</p>
    </div>
  );
}