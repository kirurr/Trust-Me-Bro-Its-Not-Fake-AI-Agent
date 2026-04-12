import { Card } from "antd";
import { useGetUsersQuery } from "../../users/userApi";
import ChatList from "./ChatList";
import ChatWindow from "./ChatWindow";

export default function ChatPage() {
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
      </div>
			<ChatWindow />
    </Card>
  );
}
