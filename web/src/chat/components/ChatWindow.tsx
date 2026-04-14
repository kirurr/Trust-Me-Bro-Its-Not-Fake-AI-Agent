import { Button, Input } from "antd";
import { useAppSelector } from "../../app/hooks";
import { selectActiveUser } from "../../users/userApi";
import { useSendMessageMutation } from "../chatApi";
import { useState } from "react";
import type { WsMessage } from "../wsMessage";

export default function ChatWindow() {
  const activeUser = useAppSelector(selectActiveUser);
  const [sendMessage, { isLoading, isError, error }] = useSendMessageMutation();
  const [message, setMessage] = useState("");

  if (activeUser === undefined) {
    return <div>no user selected</div>;
  }

  const handleSend = () => {
    if (message.trim() === "") {
      return;
    }

    const messageToSend: WsMessage = {
      userId: activeUser.user.id,
      text: message,
    };
    sendMessage(messageToSend);
    setMessage("");
  };
  return (
    <div className="w-full">
      <ul>
        {activeUser.messages.map((m) => (
          <li key={m.id}>{m.message}</li>
        ))}
      </ul>
      <div className="flex flex-row gap-2">
        <Input
					value={message}
					onChange={(e) => setMessage(e.target.value)}
				/>
        <Button loading={isLoading} onClick={handleSend} type="primary">
          send
        </Button>
        {isError && <div>Error: {error.toString()}</div>}
      </div>
    </div>
  );
}
