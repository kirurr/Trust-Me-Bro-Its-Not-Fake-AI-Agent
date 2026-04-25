import { useSendMessageMutation } from "../chatApi";
import { useState } from "react";
import clsx from "clsx";
import TextArea from "antd/es/input/TextArea";
import { Button } from "antd";

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
      <Button
        type="primary"
        size="large"
        htmlType="submit"
        disabled={isLoading}
      >
        send
      </Button>
      {isError && <div>Error: {error.toString()}</div>}
    </form>
  );
}
