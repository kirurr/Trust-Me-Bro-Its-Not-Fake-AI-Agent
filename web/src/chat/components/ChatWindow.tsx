import { useEffect, useRef } from "react";
import { useAppSelector } from "../../app/hooks";
import { selectActiveUser } from "../../users/userApi";
import ChatInputForm from "./ChatInputForm";
import type { Message } from "../../users/user";
import { formatSmartDate } from "../../utils/utils";

export default function ChatWindow() {
  const user = useAppSelector(selectActiveUser);
  const bottomRef = useRef<HTMLLIElement>(null);
  const messages = user?.messages;

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "auto" });
  }, [messages]);

  if (user === undefined) {
    return (
      <div className="w-full bg-slate-900/50 text-3xl flex flex-col items-center justify-center">
        <span>no user selected</span>
      </div>
    );
  }

  return (
    <div className="w-full h-full flex flex-col min-h-0 bg-slate-900/50 p-2 rounded-md">
      <ul className="flex-1 overflow-y-auto space-y-2 min-h-0">
        {user.messages.map((m) => {
          if (m.role === "user") {
            return <UserMessage key={m.id} message={m} />;
          } else {
            return <SystemMessage key={m.id} message={m} />;
          }
        })}
        {messages && <li ref={bottomRef} />}
      </ul>
      <ChatInputForm userId={user.user.id} />
    </div>
  );
}

function UserMessage({ message }: { message: Message }) {
  return (
    <li className="flex flex-row items-center gap-4 ml-auto w-fit">
      <span className="text-gray-400">
        {formatSmartDate(new Date(message.sentAt))} User
      </span>
      <div className="p-2 rounded-md bg-red-900 whitespace-pre-wrap">
        {message.message}
      </div>
    </li>
  );
}
function SystemMessage({ message }: { message: Message }) {
  return (
    <li className="flex flex-row items-center gap-4 w-fit">
      <div className="p-2 rounded-md bg-blue-900 w-fit whitespace-pre-wrap">
        {message.message}
      </div>
      <span className="text-gray-400">
        {formatSmartDate(new Date(message.sentAt))} You
      </span>
    </li>
  );
}
