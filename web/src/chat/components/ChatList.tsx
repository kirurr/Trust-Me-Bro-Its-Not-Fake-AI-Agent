import type { UserWithMessages } from "../../users/user";
import { useAppDispatch } from "../../app/hooks";
import { setActiveUserId } from "../chatSlice";
import { Button } from "antd";

export default function ChatList({ users }: { users: UserWithMessages[] }) {
  const dispatch = useAppDispatch();
  return (
    <div>
      <ul>
        {users.map((user) => (
          <li key={user.user.id}>
            <Button onClick={() => dispatch(setActiveUserId(user.user.id))}>
              {user.user.id}
            </Button>
          </li>
        ))}
      </ul>
    </div>
  );
}
