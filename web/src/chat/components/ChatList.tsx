import type { UserWithMessages } from "../../users/user";
import { useAppDispatch } from "../../app/hooks";
import { setActiveChat } from "../chatSlice";
import { Badge, Button } from "antd";

export default function ChatList({ users }: { users: UserWithMessages[] }) {
  const dispatch = useAppDispatch();
  return (
    <div>
      <ul>
        {users.map((user) => (
          <li key={user.user.id}>
            <Button onClick={() => dispatch(setActiveChat(user.user.id))}>
              {user.user.id}
            </Button>
            {user.hasNewMessages && <Badge count={1} />}
          </li>
        ))}
      </ul>
    </div>
  );
}
