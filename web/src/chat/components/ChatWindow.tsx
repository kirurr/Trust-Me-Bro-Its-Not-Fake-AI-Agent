import { Button, Input } from "antd";
import { useAppSelector } from "../../app/hooks";
import { selectActiveUser } from "../../users/userApi";

export default function ChatWindow() {
  const activeUser = useAppSelector(selectActiveUser);

  if (activeUser === undefined) {
    return <div>no user selected</div>;
  }
  return (
    <div className="w-full">
      <ul>
        {activeUser.messages.map((m) => (
          <li key={m.id}>{m.message}</li>
        ))}
      </ul>
      <div className="flex flex-row gap-2">
        <Input />
        <Button type="primary">send</Button>
      </div>
    </div>
  );
}
