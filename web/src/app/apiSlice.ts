import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
const backendUrl = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";

export const apiSlice = createApi({
  reducerPath: "api",
  baseQuery: fetchBaseQuery({ baseUrl: backendUrl }),
	endpoints: () => ({}),
});
