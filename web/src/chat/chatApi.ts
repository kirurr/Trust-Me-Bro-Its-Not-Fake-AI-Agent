import { apiSlice } from "../app/apiSlice";
import { store } from "../app/store";
import {
  messageSchema,
  messageToJson,
  type Message,
  type UserWithMessages,
} from "../users/user";
import { setUserHasNewMessagesThunk, userApi } from "../users/userApi";
import { closeWs, getWs } from "./ws";

export const chatApi = apiSlice.injectEndpoints({
  endpoints: (builder) => ({
    sendMessage: builder.mutation<void, { userId: string; text: string }>({
      queryFn: (message) => {
        const ws = getWs();

        return new Promise((resolve, reject) => {
          if (ws.readyState !== WebSocket.OPEN) {
            reject(new Error("WebSocket is not open"));
            return;
          }

          const data = {
            id: message.userId + "-" + Math.random(),
            userId: message.userId,
            message: message.text,
            role: "system" as const,
            sentAt: new Date().toISOString(),
          };
          const json = messageToJson(data);
          ws.send(json);

          store.dispatch(
            userApi.util.updateQueryData(
              "getUsers",
              undefined,
              (draft: UserWithMessages[]) => {
                draft
                  .find((user) => user.user.id === message.userId)
                  ?.messages.push(data);
              },
            ),
          );

          resolve({ data: undefined });
        });
      },
    }),
    subscribeToUserMessages: builder.query<Message[], void>({
      queryFn: () => ({ data: [] }),

      async onCacheEntryAdded(_, { cacheDataLoaded, cacheEntryRemoved }) {
        const ws = getWs();

        try {
          await cacheDataLoaded;
          ws.addEventListener("message", (event) => {
            const json = JSON.parse(event.data);
            const message = messageSchema.parse(json);

            store.dispatch(
              userApi.util.updateQueryData(
                "getUsers",
                undefined,
                (draft: UserWithMessages[]) => {
                  const user = draft.find(
                    (user) => user.user.id === message.userId,
                  );

                  if (!user) {
                    draft.push({
                      user: {
                        id: message.userId,
                      },
                      messages: [message],
                      hasNewMessages: true,
                    });
                    return;
                  }

                  const existingMessage = user.messages.find(
                    (m) => m.id === message.id,
                  );

                  if (existingMessage) {
                    Object.assign(existingMessage, message);
                  } else {
                    user.messages.push(message);
                  }

                  store.dispatch(
                    setUserHasNewMessagesThunk(message.userId, true),
                  );
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
