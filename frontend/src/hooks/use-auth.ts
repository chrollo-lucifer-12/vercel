import { logoutAction, refreshAction } from "@/actions/auth";
import { profileAction } from "@/actions/user";
import {
  signinMutationOptions,
  signupMutationOptions,
  TOKEN_KEY,
  USER_KEY,
  profileQueryOptions,
  SESSION_KEY,
  tokenQueryOptions,
  verifyMutationOptions,
} from "@/lib/query-options";
import {
  AccessTokenDetails,
  AuthUserDetails,
  SessionDetails,
  TokenDetails,
  User,
} from "@/lib/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export const useVerify = () => {
  return useMutation({
    ...verifyMutationOptions(),
    throwOnError: false,
    onError: (error) => {
      toast.error(error.message);
    },
  });
};

export const useSignUp = () => {
  const router = useRouter();
  return useMutation({
    ...signupMutationOptions(),
    throwOnError: false,
    onError: (error) => {
      console.log(error);
      toast.error(error.message);
    },
    onSuccess: (data) => {
      if (data.success) {
        toast.success(
          "Sign up successful check your email for verification link.",
        );
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
    onError: (error) => {
      toast.error(error.message);
    },
    onSuccess: (data) => {
      if (data.success) {
        const authDetails = data.data as AuthUserDetails;
        queryClient.setQueryData(USER_KEY, {
          email: authDetails.User.email,
          name: authDetails.User.name,
        } as User);
        queryClient.setQueryData(TOKEN_KEY, {
          access_token: authDetails.access_token,
          access_token_expires_at: authDetails.access_token_expires_at,
          session_id: authDetails.session_id,
        } as TokenDetails);
        toast.success("Login successful");
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
        TOKEN_KEY,
      ) as AccessTokenDetails & SessionDetails;
      const sessionId = authData.session_id;

      if (!sessionId) throw new Error("No session ID");

      const result = await logoutAction(sessionId);
      if (!result.success) throw new Error(result.error);
    },
    onSuccess: () => {
      queryClient.cancelQueries();
      queryClient.removeQueries();
      router.push("/");
    },
  });
};

export const useAccessToken = () => {
  const queryClient = useQueryClient();
  return useQuery({
    ...tokenQueryOptions(),
    queryFn: async () => {
      const res = await refreshAction();
      if (res.success) {
        queryClient.setQueryData(TOKEN_KEY, {
          access_token: res.access_token?.access_token,
          access_token_expires_at: res.access_token?.access_token_expires_at,
          session_id: res.access_token?.session_id,
        } as TokenDetails);
        return res.access_token;
      }
    },
  });
};

export const useProfile = () => {
  const queryClient = useQueryClient();

  return useQuery({
    ...profileQueryOptions(),
    queryFn: async () => {
      const tokenData = queryClient.getQueryData(TOKEN_KEY) as TokenDetails;
      const res = await profileAction(tokenData.access_token);
      if (!res.success) {
        throw new Error(res.error ?? "Failed to fetch profile");
      }
      return res.user;
    },
    initialData: queryClient.getQueryData(USER_KEY),
  });
};
