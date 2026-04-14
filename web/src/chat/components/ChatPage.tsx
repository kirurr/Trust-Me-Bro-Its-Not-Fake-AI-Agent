import { Card } from "antd";
import { useGetUsersQuery } from "../../users/userApi";
import ChatList from "./ChatList";
import ChatWindow from "./ChatWindow";
import { useSubscribeToUserMessagesQuery } from "../chatApi";

export default function ChatPage() {
	const {isError: isWsError, error: wsError} = useSubscribeToUserMessagesQuery();
  const { data, isLoading, isSuccess, isError, error } = useGetUsersQuery();

  return (
    <Card
			className="w-7xl mx-auto my-8 h-[calc(100vh-4rem)]"
			classNames={{
				body: "flex flex-1 gap-4 h-full",
			}}
		>
      <div className="w-full max-w-1/4">
				{isWsError && <div>Error: {wsError.toString()}</div>}
        {isLoading && <div>Loading...</div>}
        {isError && <div>Error: {error.toString()}</div>}

        {isSuccess && <ChatList users={data} />}
      </div>
			<ChatWindow />
    </Card>
  );
}
