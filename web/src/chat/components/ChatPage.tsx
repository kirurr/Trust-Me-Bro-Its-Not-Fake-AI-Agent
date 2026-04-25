import { useGetUsersQuery } from "../../users/userApi";
import ChatList from "./ChatList";
import ChatWindow from "./ChatWindow";
import { useSubscribeToUserMessagesQuery } from "../chatApi";

export default function ChatPage() {
  const { isError: isWsError, error: wsError } =
    useSubscribeToUserMessagesQuery();
  const { data, isLoading, isSuccess, isError, error } = useGetUsersQuery();

  return (
    <div className="w-full max-w-[90rem] mx-auto h-[calc(100vh-4rem)] bg-black/80 backdrop-blur-3xl p-4 flex flex-col gap-4 overflow-hidden rounded-xl">
      <h1 className="text-3xl font-bold text-center shrink-0">
        <span className="text-red-500">REMEMBER:</span> YOU ARE REAL HUMAN NOT A{" "}
        <span className="text-red-500">FAKE AI AGENT</span>
      </h1>
      <div className="flex gap-4 flex-1 min-h-0 overflow-hidden">
        <div className="w-1/4 overflow-hidden flex flex-col items-center gap-4">
          <div className="text-center">
            <h2 className="text-bold text-xl text-center pb-2">
              real humans to chat with
            </h2>
          </div>
          {isWsError && <div>WS Error: {JSON.stringify(wsError, null, 2)}</div>}
          {isLoading && <div>Loading...</div>}
          {isError && <div>Error: {JSON.stringify(error, null, 2)}</div>}

          {isSuccess && <ChatList users={data} />}
          <div className="text-center">
            <span className="text-red-500 text-bold block">warning:</span>
            <em>dont tell them you are not AI</em>
          </div>
        </div>
        <ChatWindow />
      </div>
    </div>
  );
}
