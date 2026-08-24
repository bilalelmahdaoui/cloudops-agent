import {
  type FormEvent,
  useState,
} from "react";

interface ChatInputProps {
  onSend: (message: string) => Promise<void>;
  disabled: boolean;
}

export function ChatInput({
  onSend,
  disabled,
}: ChatInputProps) {
  const [message, setMessage] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const trimmedMessage = message.trim();

    if (!trimmedMessage) {
      return;
    }

    setMessage("");
    void onSend(trimmedMessage);
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
        disabled={disabled}
      />

      <button
        type="submit"
        className="button button--primary"
        disabled={disabled || !message.trim()}
      >
        {disabled ? "Envoi..." : "Envoyer"}
      </button>
    </form>
  );
}
