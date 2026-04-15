import { useSendMessageMutation } from "../chatApi";
import { useState } from "react";
import clsx from "clsx";
import TextArea from "antd/es/input/TextArea";

export default function ChatInputForm({
  userId,
  className,
}: {
  userId: string;
  className?: string;
}) {
  const [sendMessage, { isLoading, isError, error }] = useSendMessageMutation();
  const [message, setMessage] = useState("");

  const handleSend = () => {
    if (message.trim() === "") {
      return;
    }

    sendMessage({
      userId: userId,
      text: message,
    });
    setMessage("");
  };
  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        handleSend();
      }}
      className={clsx(
        "flex flex-row gap-2 mt-2 shrink-0 items-center",
        className,
      )}
    >
      <TextArea
        value={message}
        className="scrollbar"
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Type your message here..."
        autoSize
        size="large"
      />
      <ChatButton isLoading={isLoading} />
      {isError && <div>Error: {error.toString()}</div>}
    </form>
  );
}

function ChatButton({ isLoading }: { isLoading: boolean }) {
  return (
    <button
			disabled={isLoading}
			type="submit"
			className="i-button animated"
		>
      send
    </button>
  );
}
