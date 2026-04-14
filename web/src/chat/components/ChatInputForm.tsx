import { Button, Input } from "antd";
import { useSendMessageMutation } from "../chatApi";
import { useState } from "react";
import clsx from "clsx";

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
      className={clsx("flex flex-row gap-2 mt-2 shrink-0", className)}
    >
      <Input value={message} onChange={(e) => setMessage(e.target.value)} />
      <Button loading={isLoading} type="primary">
        send
      </Button>
      {isError && <div>Error: {error.toString()}</div>}
    </form>
  );
}
