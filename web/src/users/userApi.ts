import { createSelector } from "@reduxjs/toolkit/react";
import { apiSlice } from "../app/apiSlice";
import type { AppThunk, RootState } from "../app/store";
import { userWithMessagesSchema, type UserWithMessages } from "../users/user";
import z from "zod";
import { getActiveUserId } from "../chat/chatSlice";

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

export const setUserHasNewMessagesThunk =
  (userId: string, input: boolean): AppThunk =>
  (dispatch, getState) => {
    const state = getState();
    const activeUserId = getActiveUserId(state);
    dispatch(
      userApi.util.updateQueryData("getUsers", undefined, (draft) => {
        const user = draft.find((u) => u.user.id === userId);
        if (!user) return;

        if (user.user.id === activeUserId) {
          user.hasNewMessages = false;
          return;
        }

        user.hasNewMessages = input;
      }),
    );
  };

export const { useGetUsersQuery } = userApi;
