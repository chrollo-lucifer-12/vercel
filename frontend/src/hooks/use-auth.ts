import { logoutAction } from "@/actions/auth";
import {
  AUTH_USER_KEY,
  signinMutationOptions,
  signupMutationOptions,
} from "@/lib/query-options";
import { AuthUserDetails } from "@/lib/types";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

export const useSignUp = () => {
  const router = useRouter();
  return useMutation({
    ...signupMutationOptions(),
    onSuccess: (data) => {
      if (data.success) {
        router.push("/auth/signin");
      }
    },
  });
};

export const useSignIn = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  return useMutation({
    ...signinMutationOptions(),
    onSuccess: (data) => {
      if (data.success) {
        queryClient.setQueryData(AUTH_USER_KEY, data.data as AuthUserDetails);
        router.push("/");
      }
    },
  });
};

export const useSignout = () => {
  const queryClient = useQueryClient();
  const router = useRouter();

  return useMutation({
    mutationFn: async () => {
      const authData = queryClient.getQueryData(
        AUTH_USER_KEY,
      ) as AuthUserDetails;
      const sessionId = authData.session_id;

      if (!sessionId) throw new Error("No session ID");

      const result = await logoutAction(sessionId);
      if (!result.success) throw new Error(result.error);
    },
    onSuccess: () => {
      queryClient.setQueryData(AUTH_USER_KEY, null);
      router.push("/");
    },
  });
};
