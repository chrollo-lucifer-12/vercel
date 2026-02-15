import { signinAction, signupAction } from "@/actions/auth";
import { mutationOptions, queryOptions } from "@tanstack/react-query";

const SIGNUP_KEY = "signup";
const SIGNIN_KEY = "signin";
export const USER_KEY = ["auth", "user"];
export const TOKEN_KEY = ["auth", "token"];

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
  });
};
