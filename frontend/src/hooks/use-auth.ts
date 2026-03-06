import { logoutAction } from "@/actions/auth";
import {
  signinMutationOptions,
  signupMutationOptions,
  profileQueryOptions,
  tokenQueryOptions,
  verifyMutationOptions,
} from "@/lib/query/query-options";
import { AccessTokenDetails, TokenDetails, User } from "@/lib/types";
import { useMutation, useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { getQueryClient } from "@/lib/query/query-provider";
import { profile } from "@/lib/axios/user";
import { clientEnv } from "@/lib/env/client";
import { da } from "zod/v4/locales";

const TOKEN_QUERY_KEY = ["token"];
const USER_QUERY_KEY = ["user"];

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
    onSuccess: () => {
      toast.success(
        "Sign up successful! Check your email for verification link.",
      );
      router.push("/signin");
    },
  });
};

export const useSignIn = () => {
  const router = useRouter();
  const queryClient = getQueryClient();

  return useMutation({
    ...signinMutationOptions(),
    onError: (error) => {
      toast.error(error.message);
    },
    onSuccess: async (data) => {
      if (data.success) {
        toast.success("Login successful");
        await queryClient.invalidateQueries({ queryKey: TOKEN_QUERY_KEY });
        await queryClient.invalidateQueries({ queryKey: USER_QUERY_KEY });
        router.push("/");
      }
    },
  });
};

export const useSignout = () => {
  const queryClient = getQueryClient();
  const router = useRouter();
  const { data } = useSession();
  const sessionData = data as TokenDetails;

  return useMutation({
    mutationFn: async () => {
      if (!sessionData?.session_id) throw new Error("No session ID");

      const result = await logoutAction(
        sessionData.session_id,
        sessionData.access_token,
      );
      if (!result.success) throw new Error(result.error);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: TOKEN_QUERY_KEY });
      await queryClient.invalidateQueries({ queryKey: USER_QUERY_KEY });
      router.push("/");
    },
  });
};

export const useSession = () => {
  return useQuery({
    ...tokenQueryOptions(),
    queryKey: TOKEN_QUERY_KEY,

    queryFn: async (): Promise<TokenDetails | null> => {
      const stored = sessionStorage.getItem("session_data");

      if (stored) {
        const parsed = JSON.parse(stored);
        if (parsed) {
          const expiresAt = new Date(parsed.access_token_expires_at).getTime();
          const timeLeft = expiresAt - Date.now();
          if (timeLeft > 2 * 60 * 1000) {
            return parsed;
          }
        }
      }

      const res = await fetch(`${clientEnv.NEXT_PUBLIC_BASE_URL}/api/refresh`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return null;
      }

      const text = await res.text();
      if (!text) return null;

      const data = JSON.parse(text);

      const sessionData = data.access_token as TokenDetails;

      // const sessionData: TokenDetails = {
      //   access_token: data.access_token,
      //   access_token_expires_at: data.access_token_expires_at,
      //   session_id: data.session_id,
      // };

      sessionStorage.setItem("session_data", JSON.stringify(sessionData));

      return data.success ? sessionData : null;
    },
  });
};

export const useProfile = () => {
  const { data } = useSession();
  const sessionData = data as TokenDetails;

  return useQuery({
    ...profileQueryOptions(),
    queryKey: ["profile"],
    queryFn: async (): Promise<User | null> => {
      try {
        if (!sessionData?.access_token) return null;
        const res = await profile(sessionData?.access_token!);
        return res;
      } catch (err) {
        console.error("Failed to fetch profile:", err);
        return null;
      }
    },
    enabled: !!sessionData?.access_token,
  });
};
