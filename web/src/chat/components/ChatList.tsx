import type { UserWithMessages } from "../../users/user";
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { setActiveChat } from "../chatSlice";
import { Badge } from "antd";
import { selectActiveUser } from "../../users/userApi";
import clsx from "clsx";

export default function ChatList({ users }: { users: UserWithMessages[] }) {
  const dispatch = useAppDispatch();
  const activeUser = useAppSelector(selectActiveUser);

  return (
    <div className="h-full overflow-y-auto overflow-x-hidden">
      <ul className="flex flex-col">
        {users.map((user) => {
          const isActive = user.user.id === activeUser?.user.id;
          return (
            <li
              key={user.user.id}
              onClick={() => dispatch(setActiveChat(user.user.id))}
              className={clsx(
                "p-2 flex flex-row gap-2 items-center cursor-pointer rounded-2xl hover:bg-gray-700",
                isActive && "bg-gray-800",
              )}
            >
              {user.user.id}
              {user.hasNewMessages && <Badge count={1} />}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
