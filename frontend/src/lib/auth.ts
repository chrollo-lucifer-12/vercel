import "server-only";

import { axiosInstance } from "./axios";
import { LoginInput } from "./types";
import { serverEnv } from "./env/server";

export const login = async (data: LoginInput) => {
  const { email, name, password } = data;

  await axiosInstance.post(serverEnv.LOGIN_ENDPOINT);
};
