import { logoutAction } from "@/actions/auth";
import {
  signinMutationOptions,
  signupMutationOptions,
  profileQueryOptions,
  tokenQueryOptions,
  verifyMutationOptions,
} from "@/lib/query-options";
import { TokenDetails, User } from "@/lib/types";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { getQueryClient } from "@/lib/query-provider";
import { profile } from "@/lib/axios/user";
import { clientEnv } from "@/lib/env/client";

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
    throwOnError: false,
    queryFn: async (): Promise<TokenDetails | null> => {
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

      return data.success ? data.access_token : null;
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
