import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { User } from '@/types';
import type { AuthState } from '@/redux/types';

const initialState: AuthState = {
  user: undefined,
  token: localStorage.getItem('access_token') || undefined,
  refreshToken: localStorage.getItem('refresh_token') || undefined,
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
      localStorage.setItem('access_token', action.payload);
    },
    clearToken: (state) => {
      state.token = undefined;
      localStorage.removeItem('access_token');
    },
    setRefreshToken: (state, action: PayloadAction<string>) => {
      state.refreshToken = action.payload;
      localStorage.setItem('refresh_token', action.payload);
    },
    clearAuth: (state) => {
      state.user = undefined;
      state.token = undefined;
      state.refreshToken = undefined;
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    },
  },
});

export const { setUser, clearUser, setToken, clearToken, setRefreshToken, clearAuth } = authSlice.actions;
export default authSlice.reducer;
