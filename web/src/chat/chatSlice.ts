import { createSlice, type PayloadAction } from "@reduxjs/toolkit";

interface ChatState {
  activeUserId: string | null;
}

const initialState: ChatState = {
  activeUserId: null,
};

export const chatSlice = createSlice({
  name: "chat",
  initialState,
  reducers: {
    setActiveUserId: (state, action: PayloadAction<string>) => {
      state.activeUserId = action.payload;
    },
    clearActiveUserId: (state) => {
      state.activeUserId = null;
    },
  },
  selectors: {
    getActiveUserId: (state) => state.activeUserId,
  },
});

export const { setActiveUserId, clearActiveUserId } = chatSlice.actions;
export const { getActiveUserId } = chatSlice.selectors;

export default chatSlice.reducer;
