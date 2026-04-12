import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { userSchema, type User } from "../users/user";
import z from "zod";
const backendUrl = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const apiSlice = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({ baseUrl: backendUrl }),

  endpoints: (builder) => ({
    getUsers: builder.query<User[], void>({
      query: () => "/users",

      transformResponse: (response: unknown) => {
        return z.array(userSchema).parse(response);
      },
    }),
  }),
});

export const { useGetUsersQuery } = apiSlice;
