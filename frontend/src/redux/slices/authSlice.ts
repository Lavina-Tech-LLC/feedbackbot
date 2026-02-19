import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { User } from '@/types';
import type { AuthState } from '@/redux/types';

const initialState: AuthState = {
  user: undefined,
  token: undefined,
};

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload;
    },
    clearUser: (state) => {
      state.user = undefined;
    },
    setToken: (state, action: PayloadAction<string>) => {
      state.token = action.payload;
    },
    clearToken: (state) => {
      state.token = undefined;
    },
    clearAuth: (state) => {
      state.user = undefined;
      state.token = undefined;
    },
  },
});

export const { setUser, clearUser, setToken, clearToken, clearAuth } = authSlice.actions;
export default authSlice.reducer;
