import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

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

      {message.role === "assistant" ? (
        <div className="message__content">
          <Markdown remarkPlugins={[remarkGfm]}>
            {message.content}
          </Markdown>
        </div>
      ) : (
        <p>{message.content}</p>
      )}
    </div>
  );
}
