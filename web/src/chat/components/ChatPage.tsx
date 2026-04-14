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
			className="max-w-5xl mx-auto mt-8"
			classNames={{
				body: "flex flex-row gap-4"
			}}
		>
      <div>
        {isLoading && <div>Loading...</div>}
        {isError && <div>Error: {error.toString()}</div>}

        {isSuccess && <ChatList users={data} />}
				{isWsError && <div>Error: {wsError.toString()}</div>}
      </div>
			<ChatWindow />
    </Card>
  );
}
