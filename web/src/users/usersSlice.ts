import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Message, User } from "./user";

export type UserState = {
  user: User;
  messages: Message[];
};

export interface UsersState {
  users: UserState[];
}

const initialState: UsersState = {
  users: [],
};

export const usersSlice = createSlice({
  name: "users",
  initialState,
  reducers: {
    addUser: (state, action: PayloadAction<UserState>) => {
      state.users = [...state.users, action.payload];
    },
  },
	selectors: {
		getUsers: (state) => state.users,
		getUserById: (state, id) => state.users.find((user) => user.user.id === id),
	}
});

export const { addUser } = usersSlice.actions;
export const { getUsers, getUserById } = usersSlice.selectors;

export default usersSlice.reducer;
