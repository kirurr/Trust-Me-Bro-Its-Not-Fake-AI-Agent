import { configureStore } from "@reduxjs/toolkit";
import usersSlice from "../users/usersSlice";
import { listenerMiddleware } from "./listenerMiddleware";
import { apiSlice } from "./apiSlice";

export const store = configureStore({
  reducer: {
    users: usersSlice,
    [apiSlice.reducerPath]: apiSlice.reducer,
  },
  middleware: (getDefaultMIddleware) =>
    getDefaultMIddleware()
      .prepend(listenerMiddleware.middleware)
      .concat(apiSlice.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
export type AppStore = typeof store;
