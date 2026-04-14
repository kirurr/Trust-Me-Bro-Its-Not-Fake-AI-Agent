import { apiSlice } from "../app/apiSlice";
import { store } from "../app/store";
import type { UserWithMessages } from "../users/user";
import { userApi } from "../users/userApi";
import { closeWs, getWs } from "./ws";
import {
  wsMessageSchema,
  wsMessageToUserMessage,
  type WsMessage,
} from "./wsMessage";

export const chatApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    sendMessage: builder.mutation<void, WsMessage>({
      queryFn: (message) => {
        const ws = getWs();

        return new Promise((resolve, reject) => {
          if (ws.readyState !== WebSocket.OPEN) {
            reject(new Error("WebSocket is not open"));
            return;
          }

          const data = JSON.stringify({
            user_id: message.userId,
            text: message.text,
          });
          ws.send(data);
          store.dispatch(
            userApi.util.updateQueryData(
              "getUsers",
              undefined,
              (draft: UserWithMessages[]) => {
                draft
                  .find((user) => user.user.id === message.userId)
                  ?.messages.push(wsMessageToUserMessage(message, "system"));
              },
            ),
          );

          resolve({ data: undefined });
        });
      },
    }),
    subscribeToUserMessages: builder.query<WsMessage[], void>({
      queryFn: () => ({ data: [] }),

      async onCacheEntryAdded(_, { cacheDataLoaded, cacheEntryRemoved }) {
        const ws = getWs();

        try {
          await cacheDataLoaded;

          ws.addEventListener("open", () => {
            console.log("WebSocket connection established");
          });

          ws.addEventListener("close", () => {
            console.log("WebSocket connection closed");
          });

          ws.addEventListener("message", (event) => {
            const message = JSON.parse(event.data);
            const parsedWsMessage = wsMessageSchema.parse(message);
            const userMessage = wsMessageToUserMessage(parsedWsMessage, "user");

            store.dispatch(
              userApi.util.updateQueryData(
                "getUsers",
                undefined,
                (draft: UserWithMessages[]) => {
                  const user = draft.find(
                    (user) => user.user.id === userMessage.userId,
                  );

                  if (!user) {
                    draft.push({
                      user: {
                        id: userMessage.userId,
                      },
                      messages: [userMessage],
                    });
                    return;
                  }

                  user.messages.push(userMessage);
                },
              ),
            );
          });
        } catch (e) {
          console.error(e);
        }

        await cacheEntryRemoved;
        closeWs();
      },
    }),
  }),
});

export const { useSubscribeToUserMessagesQuery, useSendMessageMutation } =
  chatApi;
