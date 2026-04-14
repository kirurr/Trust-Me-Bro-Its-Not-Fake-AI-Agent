import { configureStore } from "@reduxjs/toolkit";
import { listenerMiddleware } from "./listenerMiddleware";
import { apiSlice } from "./apiSlice";
import chatSlice from "../chat/chatSlice";

export const store = configureStore({
  reducer: {
    chat: chatSlice,
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

export type AppThunk<ReturnType = void> = (
  dispatch: AppDispatch,
  getState: () => RootState,
) => ReturnType;
