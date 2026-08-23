import {
  type FormEvent,
  useState,
} from "react";

interface ChatInputProps {
  onSend: (message: string) => void;
}

export function ChatInput({ onSend }: ChatInputProps) {
  const [message, setMessage] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedMessage = message.trim();

    if (!trimmedMessage) {
      return;
    }

    onSend(trimmedMessage);
    setMessage("");
  };

  return (
    <form
      className="chat-input"
      onSubmit={handleSubmit}
    >
      <input
        value={message}
        onChange={(event) => setMessage(event.target.value)}
        placeholder="Demandez quelque chose à CloudOps Agent..."
        aria-label="Message"
      />

      <button
        type="submit"
        className="button button--primary"
        disabled={!message.trim()}
      >
        Envoyer
      </button>
    </form>
  );
}