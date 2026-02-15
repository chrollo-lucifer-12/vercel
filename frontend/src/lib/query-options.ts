import { signinAction, signupAction } from "@/actions/auth";
import { mutationOptions } from "@tanstack/react-query";

const SIGNUP_KEY = "signup";
const SIGNIN_KEY = "signin";
export const AUTH_USER_KEY = ["auth", "user"];

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
