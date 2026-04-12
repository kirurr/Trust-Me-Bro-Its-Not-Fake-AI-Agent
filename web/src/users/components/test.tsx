import { Button } from "antd";
import { useAppDispatch, useAppSelector } from "../../app/hooks";
import { addUser, getUsers } from "../usersSlice";
import { apiSlice, useGetUsersQuery } from "../../app/apiSlice";

export default function Test() {
  const dispatch = useAppDispatch();
  const users = useAppSelector(getUsers);

  const { data, isLoading, isSuccess, isError, error } = useGetUsersQuery();

  const addUserCallback = () => {
    dispatch(addUser({ user: { id: "1" + Math.random() }, messages: [] }));
  };

  const apiCallback = () => {
    dispatch(
      apiSlice.util.updateQueryData("getUsers", undefined, (draft) => {
        draft.push({ id: "1" + Math.random() });
      }),
    );
  };

  return (
    <div>
      <Button type="primary" onClick={addUserCallback}>
        Add User
      </Button>
      <Button type="primary" onClick={apiCallback}>
        Add User API
      </Button>

      <ul>
        {users.map((user) => (
          <li key={user.user.id}>{user.user.id}</li>
        ))}
      </ul>

      {isLoading && <div>Loading...</div>}
      {isError && <div>Error: {error.toString()}</div>}

      {isSuccess && <pre>{JSON.stringify(data, null, 2)}</pre>}
    </div>
  );
}
