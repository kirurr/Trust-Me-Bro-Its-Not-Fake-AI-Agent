import { createSelector } from "@reduxjs/toolkit/react";
import { apiSlice } from "../app/apiSlice";
import type { RootState } from "../app/store";
import { userWithMessagesSchema, type UserWithMessages } from "../users/user";
import z from "zod";

export const userApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    getUsers: builder.query<UserWithMessages[], void>({
      query: () => "/users",

      transformResponse: (response: unknown) => {
        return z.array(userWithMessagesSchema).parse(response);
      },
    }),
  }),
});

const selectUsersData = userApi.endpoints.getUsers.select(undefined);

export const selectActiveUser = createSelector(
  (state: RootState) => selectUsersData(state).data,
  (state: RootState) => state.chat.activeUserId,
  (users, activeUserId) => users?.find((user) => user.user.id === activeUserId),
);

export const { useGetUsersQuery } = userApi;
