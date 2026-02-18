import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import type { User } from '@/types';
import type { AuthState } from '@/redux/types';

const initialState: AuthState = {
  user: undefined,
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
  },
});

export const { setUser, clearUser } = authSlice.actions;
export default authSlice.reducer;
