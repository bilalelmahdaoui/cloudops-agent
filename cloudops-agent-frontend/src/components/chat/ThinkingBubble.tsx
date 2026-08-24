export function ThinkingBubble() {
  return (
    <div
      className="message message--assistant message--thinking"
      role="status"
      aria-label="CloudOps Agent prépare une réponse"
    >
      <span className="message__author">CloudOps Agent</span>

      <span className="thinking-dots" aria-hidden="true">
        <span />
        <span />
        <span />
      </span>
    </div>
  );
}
