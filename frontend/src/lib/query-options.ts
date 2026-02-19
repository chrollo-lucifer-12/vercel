import { signinAction, signupAction, verifyAction } from "@/actions/auth";
import { mutationOptions, queryOptions } from "@tanstack/react-query";

const SIGNUP_KEY = "signup";
const SIGNIN_KEY = "signin";
export const USER_KEY = ["auth", "user"];
export const TOKEN_KEY = ["auth", "token"];
export const SESSION_KEY = ["session", "refresh"];
const VERIFY_KEY = ["verify"];

export const verifyMutationOptions = () => {
  return mutationOptions({
    mutationKey: VERIFY_KEY,
    mutationFn: verifyAction,
  });
};

export const signupMutationOptions = () => {
  return mutationOptions({
    mutationKey: [SIGNUP_KEY],
    mutationFn: signupAction,
  });
};

export const signinMutationOptions = () => {
  return mutationOptions({
    mutationKey: [SIGNIN_KEY],
    mutationFn: signinAction,
  });
};

export const profileQueryOptions = () => {
  return queryOptions({
    queryKey: USER_KEY,
    staleTime: Infinity,
    refetchOnWindowFocus: false,
  });
};

export const tokenQueryOptions = () => {
  return queryOptions({
    queryKey: TOKEN_KEY,
    refetchInterval: 1000 * 60 * 1,
    staleTime: 1000 * 60 * 12,
    refetchOnWindowFocus: false,
    refetchIntervalInBackground: true,
  });
};
